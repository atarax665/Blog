CREATE OR REPLACE FUNCTION fn_new_bets() RETURNS TRIGGER AS 
$$
BEGIN
    PERFORM pg_notify(
        'new_bet',
        to_json(NEW)::TEXT
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER new_bet
AFTER INSERT ON bets
FOR EACH ROW EXECUTE PROCEDURE fn_new_bets();