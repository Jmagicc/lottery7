SELECT
    l1.draw_date AS draw_date_1,
    l1.num1 AS num1_1,
    l1.num2 AS num2_1,
    l2.draw_date AS draw_date_2,
    l2.num1 AS num1_2,
    l2.num2 AS num2_2,
    GREATEST(CAST(l1.draw_no AS UNSIGNED), CAST(l2.draw_no AS UNSIGNED)) -
    LEAST(CAST(l1.draw_no AS UNSIGNED), CAST(l2.draw_no AS UNSIGNED)) AS draw_gap
FROM
    lottery_results l1
JOIN
    lottery_results l2
ON
    l1.num1 = l2.num1 AND l1.num2 = l2.num2 AND l1.draw_no < l2.draw_no
and   l1.num1 = 0 and l1.num2 = 0
ORDER BY
    draw_gap ASC,  -- 先按 draw_gap 升序排列
    draw_date_2 DESC  -- 再按 draw_date_2 降序排列
LIMIT 1000;
