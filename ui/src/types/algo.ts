export const shipGenerationPerRound = (
  industry: number,
  manufacturing: number,
) => industry * (manufacturing + 5);

export const shipGenerationPerTick = (
  industry: number,
  manufacturing: number,
  productionRate: number,
) => shipGenerationPerRound(industry, manufacturing) / productionRate;

interface PositionedThing {
  x: string;
  y: string;
}

export const distanceBetween = (
  a: PositionedThing,
  b: PositionedThing,
) => {
  // pythagoras
  const dX = Math.abs(parseFloat(a.x) - parseFloat(b.x));
  const dY = Math.abs(parseFloat(a.y) - parseFloat(b.y));
  return Math.sqrt(
    (dX * dX) + (dY * dY),
  );
};

export const pointsNeededForTechLevel = (
  targetLevel: number,
) => 144 * (targetLevel - 1);

export const guessBattle = (
  attackerShips: number,
  attackerWeaponsLevel: number,
  defenderShips: number,
  defenderWeaponsLevel: number,
) => {
  const defenderWeaponsLevelWithBonus = defenderWeaponsLevel + 1; // defenders bonus

  const ticksToKillAttacker = attackerShips / defenderWeaponsLevelWithBonus; // 14.75
  const ticksToKillDefender = defenderShips / attackerWeaponsLevel; // 0

  const lowestTicks = Math.ceil(Math.min(ticksToKillAttacker, ticksToKillDefender)); // 0

  const attackerShipsRemaining = Math.max(
    0,
    attackerShips - (lowestTicks * defenderWeaponsLevelWithBonus), // 59
  );
  let defenderShipsRemaining = Math.max(
    0,
    defenderShips - (lowestTicks * attackerWeaponsLevel), // 0
  );

  if (attackerShipsRemaining === 0 && defenderShipsRemaining === 0) {
    // defender wins
    defenderShipsRemaining += attackerWeaponsLevel;
  }

  const attackerWins = attackerShipsRemaining > 0;
  const defenderWins = defenderShipsRemaining > 0;

  if (attackerWins === defenderWins) {
    throw new Error('encountered draw but this should be impossible');
  }

  // the rounding could be incorrect, not 100% sure
  const defenderShipsNeeded = Math.floor(Math.max(
    0,
    (ticksToKillAttacker * attackerWeaponsLevel) - defenderShips,
  ));

  return {
    attackerWins,
    defenderWins,

    attackerShipsRemaining,
    defenderShipsRemaining,

    lowestTicks,
    defenderShipsNeeded,
  };
};
