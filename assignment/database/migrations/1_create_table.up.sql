
-- Create Bill table
CREATE TABLE Bills (
    id UUID PRIMARY KEY,
    customer_id VARCHAR(255) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(255) NOT NULL,
    total_charges DECIMAL(10, 2) NOT NULL
);

CREATE INDEX idx_bills_id ON Bills(id);
-- Create LineItem table
CREATE TABLE LineItems (
    id UUID PRIMARY KEY,
    bill_id UUID NOT NULL,
    description TEXT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    currency VARCHAR(3) NOT NULL,
    metadata TEXT,
    FOREIGN KEY (bill_id) REFERENCES Bills(id) ON DELETE CASCADE
);
CREATE INDEX idx_lineitems_id ON LineItems(id);
-- Create index on bill_id in LineItems table for faster queries
CREATE INDEX idx_lineitems_bill_id ON LineItems(bill_id);