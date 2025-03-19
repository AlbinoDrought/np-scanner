// ==UserScript==
// @name     Embed NP-Scanner data into the Neptune's Pride map
// @homepageURL https://github.com/AlbinoDrought/np-scanner
// @version  1
// @grant    none
// @include https://np.ironhelmet.com/game/*
// @license AGPL-3.0
// ==/UserScript==

async function loadData(NeptunesPride, match) {
  const resp = await fetch(`${match.url}/api/matches/${match.gameId}/merged-snapshot?access_code=${match.code}`);
  if (resp.status === 401) {
    const err = new Error(`Bad Response ${resp.status}: ${resp}`);
    console.error('[np-scanner]', err, resp, await resp.text());
    throw err;
  }
  
  const apiResponse = await resp.json();

  // import better star data
  const starsToRename = [];
  Object.keys(apiResponse.scanning_data.stars).forEach((uid) => {
    var gameStateStar = NeptunesPride.universe.galaxy.stars[uid];
    if (gameStateStar.v === '1' || gameStateStar.v === 1) {
      // already visible normally
      return;
    }

    var apiResponseStar = apiResponse.scanning_data.stars[uid];
    if (!apiResponseStar || (apiResponseStar.v !== '1' && apiResponseStar.v !== 1)) {
      // not visible in merged api response
      return;
    }

    gameStateStar.v = apiResponseStar.v; // uses strings
    gameStateStar.st = apiResponseStar.st; // strength
    // I think this field should be the star power + orbiting fleet power.
    // it gets recalculated when we call NeptunesPride.universe.addGalaxy
    gameStateStar.totalDefenses = gameStateStar.st;
    gameStateStar.e = apiResponseStar.e; // econ
    gameStateStar.i = apiResponseStar.i; // industry
    gameStateStar.s = apiResponseStar.s; // science
    gameStateStar.c = apiResponseStar.c; // no idea
    gameStateStar.nr = apiResponseStar.nr; // natural resources
    gameStateStar.r = apiResponseStar.r; // resources
    gameStateStar.ga = apiResponseStar.ga; // warp gate presence
    starsToRename.push(uid);

    // cloneInto prevents data access errors
    NeptunesPride.universe.galaxy.stars[uid] = cloneInto(
      gameStateStar,
      window,
    );
  });

  // import fleets
  const fleetsToRename = [];
  Object.keys(apiResponse.scanning_data.fleets).forEach((uid) => {
    var gameStateFleet = NeptunesPride.universe.galaxy.fleets[uid];
    if (gameStateFleet) {
      // already visible normally
      return;
    }

    var apiResponseFleet = apiResponse.scanning_data.fleets[uid];
    if (!apiResponseFleet) {
      // not visible merged in api response
      return;
    }

    // if this ran multiple times, values would get wrapped multiple times
    apiResponseFleet.orders = apiResponseFleet.o;
    apiResponseFleet.player = NeptunesPride.universe.galaxy.players[apiResponseFleet.puid];
    // not sure if this is correct, but makes some sense
    apiResponseFleet.path = apiResponseFleet.orders
      .map(o => o[1]) // grab star ID
      .map(sid => NeptunesPride.universe.galaxy.stars[sid]); // map to star
    fleetsToRename.push(uid);

    // cloneInto prevents data access errors
    NeptunesPride.universe.galaxy.fleets[uid] = cloneInto(
      apiResponseFleet,
      window,
    );
    // apiResponseFleet is not close to a 1:1 match with real data, but it renders
  });

  // make game reload and hopefully fix the data
  NeptunesPride.universe.addGalaxy(NeptunesPride.universe.galaxy);

  // make the modified stars and fleets stand out, so we know the data is from np-scanner
  starsToRename.forEach((uid) => {
    NeptunesPride.universe.galaxy.stars[uid].n = `[${NeptunesPride.universe.galaxy.stars[uid].n}]`;
    NeptunesPride.universe.galaxy.stars[uid].st = `[${NeptunesPride.universe.galaxy.stars[uid].st}]`;
  });

  fleetsToRename.forEach((uid) => {
    NeptunesPride.universe.galaxy.fleets[uid].n = `[${NeptunesPride.universe.galaxy.fleets[uid].n}]`;
    NeptunesPride.universe.galaxy.fleets[uid].st = `[${NeptunesPride.universe.galaxy.fleets[uid].st}]`;
  });

  // trigger redraw
  NeptunesPride.npui.map.createSpritesStars();
  NeptunesPride.npui.map.createSpritesFleets();
  NeptunesPride.npui.map.draw();
  console.log('[np-scanner]', 'clean init finished');
}

async function boot(window) {
  if (!window.NeptunesPride) {
    throw new Error('Entire window.NeptunesPride missing');
  }

  if (!window.NeptunesPride.gameId) {
    throw new Error('window.NeptunesPride.gameId missing');
  }

  var state = JSON.parse(localStorage.getItem('np-scanner')) || {
    lastUrl: 'https://localhost',
    lastCode: '',
    matches: {},
  };
  
  var match = state.matches[window.NeptunesPride.gameId] || {
    prompted: false,
    disabled: false,
    askedForCreds: false,
    url: '',
    code: '',
  };
  match.gameId = window.NeptunesPride.gameId;
  
  if (window.location.hash === '#wipe-np-scanner') {
    // if user entered wrong code, allow reprompt at #wipe-np-scanner
    match.prompted = false;
    match.askedForCreds = false;
  }
  
  if (!match.prompted) {
    match.prompted = true;
    match.enabled = window.confirm('Enable NP-Scanner for this match?');
  
    state.matches[window.NeptunesPride.gameId] = match;
    localStorage.setItem('np-scanner', JSON.stringify(state));
  }
  
  if (!match.enabled) {
    console.log('[np-scanner]', 'match disabled');
    return;
  }
  console.log('[np-scanner]', 'match enabled');
  
  if (!match.askedForCreds) {
    match.askedForCreds = true;
    
    match.url = window.prompt('Enter NP-Scanner url', state.lastUrl);
    state.lastUrl = match.url;
  
    match.code = window.prompt('Enter NP-Scanner access code', state.lastCode);
    state.lastCode = match.lastCode;
  
    state.matches[window.NeptunesPride.gameId] = match;
    localStorage.setItem('np-scanner', JSON.stringify(state));
  }
  
  await loadData(window.NeptunesPride, match);

  // when universe data is reloaded, also reload np-scanner data
  // (you can test this by clicking the "Credits: $123 / Production: 24h" bar at the top)
  // window.NeptunesPride.np.on("order:full_universe", () => loadData(window.NeptunesPride, match));
  // (this doesn't work - I think because of unsafeWindow - polling for changes instead)
  let lastLoad = window.NeptunesPride.universe.now;
  setInterval(async () => {
    if (window.NeptunesPride.universe.now === lastLoad) {
      return;
    }

    await loadData(window.NeptunesPride, match);
    lastLoad = window.NeptunesPride.universe.now;
  }, 1000);
}

// wait for NeptunesPride global object to appear
// we need to use unsafeWindow because we inject data
let attempts = 0;
const bootHandle = setInterval (() => {
  if (typeof unsafeWindow.NeptunesPride === "object") {
    clearInterval(bootHandle);
    boot(unsafeWindow);
    return;
  }

  attempts++;
  if (attempts > 100) {
    clearInterval(bootHandle);
    console.error('[np-scanner]', 'Unable to find NeptunesPride');
  }
}, 10);
