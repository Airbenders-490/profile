-- creates tags table
CREATE TABLE tags (
    name text Primary Key,
    positive bool
);

-- sample insert queries in tags
INSERT INTO tags (name, positive) VALUES ('hardworking', true);
INSERT INTO tags (name, positive) VALUES ('slacker', false);
INSERT INTO tags (name, positive) VALUES ('leader', true);
INSERT INTO tags (name, positive) VALUES ('friendly', true);