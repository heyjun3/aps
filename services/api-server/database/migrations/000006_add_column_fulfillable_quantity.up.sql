ALTER TABLE inventories ADD COLUMN fulfillable_quantity bigint;
UPDATE inventories SET fulfillable_quantity = 0;
ALTER TABLE inventories ALTER COLUMN fulfillable_quantity SET NOT NULL;