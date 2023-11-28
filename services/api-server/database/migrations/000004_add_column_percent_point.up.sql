ALTER TABLE current_prices ADD COLUMN percent_point bigint;
ALTER TABLE lowest_prices ADD COLUMN percent_point bigint;

UPDATE current_prices SET percent_point = 0;
UPDATE lowest_prices SET percent_point = 0;

ALTER TABLE current_prices ALTER COLUMN percent_point SET NOT NULL;
ALTER TABLE lowest_prices ALTER COLUMN percent_point SET NOT NULL;
