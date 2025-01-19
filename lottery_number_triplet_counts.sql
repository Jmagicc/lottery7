create or replace
algorithm = UNDEFINED view `lottery_number_triplet_counts` as
select
    concat(`lottery_results`.`num1`, `lottery_results`.`num2`, `lottery_results`.`num3`) as `num_triplet`,
    count(0) as `triplet_count`
from
    `lottery_results`
group by
    `num_triplet`;
