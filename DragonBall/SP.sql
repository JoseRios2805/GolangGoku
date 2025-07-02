CREATE OR REPLACE FUNCTION save_character(
    p_ID_character INTEGER,
    p_name TEXT,
    p_description TEXT
)
RETURNS INTEGER AS $$
DECLARE
    existing_id INTEGER;
BEGIN
    SELECT id INTO existing_id
    FROM Characters
    WHERE ID_character = p_ID_character;

    IF existing_id IS NOT NULL THEN
        RETURN existing_id;
    ELSE
        INSERT INTO Characters (ID_character, name, description)
        VALUES (p_ID_character, p_name, p_description)
        RETURNING id INTO existing_id;

        RETURN existing_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_character_by_name(p_name TEXT)
RETURNS TABLE (
    id INTEGER,
    ID_character INTEGER,
    name TEXT,
    description TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        c.id,
        c.ID_character,
        c.name,
        c.description
    FROM Characters c
    WHERE c.name ILIKE '%' || p_name || '%'
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;