SELECT
  Challenge, Attempts, Successes,
  pvalue,
  RANK() OVER (ORDER BY pvalue) Difficulty
FROM binomial_test(
  SELECT
    id,
    total,
    observed,
    sum_observed/sum_total probability
  FROM (
    SELECT
      Challenge id,
      Attempts total,
      Successes observed,
      SUM(Attempts) OVER () sum_total,
      SUM(Successes) OVER () sum_observed
    FROM tmp.man_v_food
    WHERE Attempts > 0
  )
) TestResults
JOIN (
  SELECT
    Challenge,
    Attempts,
    Successes
  FROM tmp.man_v_food
) ChallengeData
ON TestResults.id == ChallengeData.Challenge
;
