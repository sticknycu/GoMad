-- Create table products
CREATE TABLE IF NOT EXISTS products(
    id varchar primary key,
    name char(64) not null,
    manufacturer char(64),
    price int,
    stock int,
    tags char(64)[]
    );

insert into products(id, name, manufacturer, price, stock, tags) VALUES ('sha256', 'Product1', 'Manufacturer1', 15, 3, ARRAY['TAG1', 'TAG2'])

drop table products;

-- Example insert
insert into products(id, name, tags) values('sha256', 'Lapte', ARRAY['lactate', 'uht']);

-- Example selects
select * from products;
select id, name, tags from products ;
select id, name, from products where name='Lapte';

-- Example update
update products set price = 15000 WHERE name = 'Lapte';

-- Example delete
delete from products where name = 'Lapte';
