-- ADding Foreing Keys
ALTER TABLE workshops
ADD CONSTRAINT fk_category
FOREIGN KEY (category_id)
REFERENCES training_categories(id)
ON DELETE RESTRICT;

ALTER TABLE workshops
ADD CONSTRAINT fk_type
FOREIGN KEY (type_id)
REFERENCES training_types(id)
ON DELETE RESTRICT;