CREATE TABLE IF NOT EXISTS expenses (
    id            TEXT            PRIMARY KEY,
    supplier      TEXT            NOT NULL,
    issued_at     TIMESTAMPTZ     NOT NULL,
    description   TEXT,
    qty           INTEGER         NOT NULL,
    cost_per_unit NUMERIC(12, 2)  NOT NULL,
    category      TEXT            NOT NULL,
    created_at    TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_expenses_supplier ON expenses (supplier);
CREATE INDEX idx_expenses_issued_at ON expenses (issued_at);

DROP FUNCTION IF EXISTS create_expense;
CREATE FUNCTION create_expense(
    expense JSON
) RETURNS JSON AS $$
DECLARE
    r       expenses;
    created expenses;
BEGIN
    r := json_populate_record(NULL::expenses, expense);

    INSERT INTO expenses (id, supplier, issued_at, description, qty, cost_per_unit, category)
    VALUES (r.id, r.supplier, r.issued_at, r.description, r.qty, r.cost_per_unit, r.category)
    RETURNING * INTO created;

    RETURN row_to_json(created);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

DROP FUNCTION IF EXISTS update_expense;
CREATE FUNCTION update_expense(
    expense JSON
) RETURNS SETOF JSON AS $$
DECLARE
    r expenses;
BEGIN
    r := json_populate_record(NULL::expenses, expense);

    RETURN QUERY
    WITH updated AS (
        UPDATE expenses AS e
            SET supplier      = r.supplier,
                issued_at     = r.issued_at,
                description   = r.description,
                qty           = r.qty,
                cost_per_unit = r.cost_per_unit,
                category      = r.category,
                updated_at    = NOW()
            WHERE e.id = r.id
            RETURNING e.*
    )
    SELECT row_to_json(updated) FROM updated;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

DROP FUNCTION IF EXISTS delete_expense;
CREATE FUNCTION delete_expense(
    expense_id TEXT
) RETURNS VOID AS $$
BEGIN
    DELETE FROM expenses WHERE id = expense_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

DROP FUNCTION IF EXISTS list_expenses;
CREATE FUNCTION list_expenses(
    query_filter_text TEXT DEFAULT NULL,
    query_supplier TEXT DEFAULT NULL,
    query_from TIMESTAMPTZ DEFAULT NULL,
    query_to TIMESTAMPTZ DEFAULT NULL,
    query_limit INTEGER DEFAULT NULL,
    query_offset INTEGER DEFAULT 0
) RETURNS SETOF JSON AS $$
BEGIN
    RETURN QUERY
    SELECT row_to_json(e)
    FROM expenses e
    WHERE (query_filter_text IS NULL OR e.description ILIKE '%' || query_filter_text || '%')
      AND (query_supplier IS NULL OR e.supplier ILIKE '%' || query_supplier || '%')
      AND (query_from IS NULL OR e.issued_at >= query_from)
      AND (query_to IS NULL OR e.issued_at <= query_to)
    ORDER BY e.issued_at DESC
    LIMIT query_limit OFFSET query_offset;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

DROP FUNCTION IF EXISTS count_expenses;
CREATE FUNCTION count_expenses(
    query_filter_text TEXT DEFAULT NULL,
    query_supplier TEXT DEFAULT NULL,
    query_from TIMESTAMPTZ DEFAULT NULL,
    query_to TIMESTAMPTZ DEFAULT NULL
) RETURNS BIGINT AS $$
BEGIN
    RETURN (
        SELECT COUNT(*)
        FROM expenses e
        WHERE (query_filter_text IS NULL OR e.description ILIKE '%' || query_filter_text || '%')
          AND (query_supplier IS NULL OR e.supplier ILIKE '%' || query_supplier || '%')
          AND (query_from IS NULL OR e.issued_at >= query_from)
          AND (query_to IS NULL OR e.issued_at <= query_to)
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

DROP FUNCTION IF EXISTS find_expense;
CREATE FUNCTION find_expense(
    expense_id TEXT
) RETURNS SETOF JSON AS $$
BEGIN
    RETURN QUERY
    SELECT row_to_json(e)
    FROM expenses e
    WHERE e.id = expense_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;