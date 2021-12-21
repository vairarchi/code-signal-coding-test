CREATE PROCEDURE solution()
BEGIN
	/* Write your SQL here. Terminate each statement with a semicolon. */
    WITH RECURSIVE
    ce AS ( SELECT origin,  destination, 0 stops, cost                   
            FROM flights
        UNION ALL
            SELECT ce.origin, flights.destination, ce.stops + 1, ce.cost + flights.cost                   
            FROM ce
            JOIN flights ON ce.destination = flights.origin  AND ce.stops < 2 ),
    ce2 AS (SELECT *, RANK() OVER (PARTITION BY origin, destination  ORDER BY cost, stops) rk
            FROM ce)
    SELECT origin, destination, stops, cost as total_cost
    FROM ce2
    WHERE rk = 1
    ORDER BY 1,2,3;
END