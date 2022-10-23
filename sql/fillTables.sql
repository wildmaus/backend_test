insert into users (id, balance, reserved)
values (1, 100, 3);
insert into users (id, balance, reserved)
values (2, 400, 34);


insert into details (orderId, serviceId, status)
values (2, 3, true);
insert into details (orderId, serviceId, status)
values (4, 4, false);

insert into transactions (fromId, toId, amount, date, type, detailsId)
values (1, 1, 100, '2022-01-22', 0, 1);