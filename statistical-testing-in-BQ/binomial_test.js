/*
 * Binomial test for BigQuery
 *
 * Description
 *
 *    Performs an exact test of a simple null hypothesis that the probability of
 *    success in a Bernoulli experiment is `p` with an alternative hypothesis
 *    that the probability is less than `p`.
 *
 * Arguments
 *
 *    k   Number of successes.
 *    n   Number of trials.
 *    p   Probability of success under the null hypothesis.
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
function binomial_test(k, n, p){
  if(k < 0 || k > n || n <= 0 || p < 0 || p > 1) return NaN;
  // i = 0 term
  var logcoef = 0;
  var pvalue = Math.pow(Math.E, n*Math.log(1-p)); // Math.exp is not available
  // i > 0 terms
  for(var i = 1; i <= k; i++) {
    logcoef = logcoef + Math.log(n-i+1) - Math.log(i);
    pvalue = pvalue + Math.pow(Math.E, logcoef + i*Math.log(p) + (n-i)*Math.log(1-p));
  }
  return pvalue;
}

// Function registration
bigquery.defineFunction(
  // Name used to call the function from SQL
  'binomial_test',
  // Input column names
  [
    'id',
    'observed',
    'total',
    'probability'
  ],
  // JSON representation of the output schema
  [
    { name: 'id', type: 'string' },
    { name: 'pvalue', type: 'float' }
  ],
  // Function definition
  function(row, emit) {
    emit({
      id: row.id,
      pvalue: binomial_test(row.observed, row.total, row.probability)
    })
  }
);
