/*
 * Bayesian credible intervals for Bernoulli trial type data
 *
 * Description
 *
 *    ...
 *
 * Arguments
 *
 *    k   Number of successes.
 *    n   Number of trials.
 *    q   Quantile, e.g. 0.5 for the median and 0.75 for the upper quartile.
 *
 * Details
 *
 *    The calculation is performed as a cumulative sum over the binomial
 *    distribution. All calculations are done in logarithmic space since the
 *    factors of each term are often very large, causing variable overflow and
 *    NaN as a result.
 *
 * Example
 *
 *    SELECT
 *      id,
 *      pvalue
 *    FROM
 *      binomial_test(
 *        SELECT
 *          *
 *        FROM
 *          (SELECT "test1" AS id,   100 AS total,   10 AS observed,    3 AS expected),
 *          (SELECT "test2" AS id,  1775 AS total,    4 AS observed,    7 AS expected),
 *          (SELECT "test3" AS id, 10000 AS total, 9998 AS observed, 9999 AS expected)
 *      )
 *
 * References
 *
 *     https://en.wikipedia.org/wiki/Binomial_distribution
 *     https://cloud.google.com/bigquery/user-defined-functions
 *
 * Author
 *
 *     Christofer Backlin, https://github.com/backlin
 */
function bayesian_ci(k, n, q){
  if(k < 0 || k > n || n <= 0 || q < 0 || q > 1) return NaN;

  var pq = new Array(401);
  var dq = 1/(pq.length-1);

  // Calculate the unnormalised log posterior probability log(p(q|k,n))
  var pq_max = 0;
  for(var i = 0; i < pq.length; i++){
    pq[i] = k*Math.log(i*dq) + (n-k)*Math.log(1 - i*dq);
    if(i == 0 || pq[i] > pq_max) pq_max = pq[i];
  }
  // Treat edge cases
  if(k == 0){
    pq[0] = 2*pq[1] - pq[2];
    pq_max = pq[0];
  } else if(k == n){
    pq[pq.length - 1] = 2*pq[pq.length-2] - pq[pq.length-3];
    pq_max = pq[pq.length-1];
  }

  // Scale up posterior to avoid dropping below the limit for computational
  // representation and transform back to linear space (still not normalised).
  var pq_sum = 0;
  for(var i = 0; i < pq.length; i++){
    pq[i] = Math.pow(Math.E, pq[i] - pq_max);
    pq_sum += pq[i];
  }

  // Normalise and find the quantile
  for(var i = 1; i < pq.length; i++){
    var mass = (pq[i-1] + pq[i]) / 2 / pq_sum;
    if(q > mass){
      q -= mass;
    } else {
      return ((i - 1) + q / mass) / pq.length;
    }
  }
  return 1;
}

// Function registration
bigquery.defineFunction(
  // Name used to call the function from SQL
  'bayesian_ci',
  // Input column names
  [
    'id',
    'observed',
    'total'
  ],
  // JSON representation of the output schema
  [
    { name: 'id', type: 'string' },
    { name: 'q_lower', type: 'float' },
    { name: 'q_median', type: 'float' },
    { name: 'q_upper', type: 'float' }
  ],
  // Function definition
  function(row, emit) {
    emit({
      id: row.id,
      q_lower: bayesian_ci(row.observed, row.total, .025),
      q_median: bayesian_ci(row.observed, row.total, .5),
      q_upper: bayesian_ci(row.observed, row.total, .975)
    })
  }
);
