CREATE OR REPLACE FUNCTION before_insert_records()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.id < 0 THEN
        SELECT COALESCE(MAX(id) + 1, 1) INTO NEW.id FROM records;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_records_trigger
    BEFORE INSERT ON records
    FOR EACH ROW
EXECUTE FUNCTION before_insert_records();
