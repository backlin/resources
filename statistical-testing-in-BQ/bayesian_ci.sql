SELECT
  Challenge, Attempts, Successes,
  q_lower, q_median, q_upper
FROM (
  SELECT
    id, q_lower, q_median, q_upper
  FROM bayesian_ci(
    SELECT
      Challenge id,
      Attempts total,
      Successes observed
    FROM tmp.man_v_food
    WHERE
      Attempts IS NOT NULL
      AND Attempts > 0
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
ORDER BY q_median
;
