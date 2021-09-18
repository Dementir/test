create table region (
	region_id uuid not null primary key,
	name text not null
);
comment on table region is 'Справочник регионов';

create table channel (
	channel_id uuid not null primary key,
	name text not null
);
comment on table channel is 'Справочник каналов';

create table status (
	status_id uuid not null primary key,
	name text not null
);
comment on table status is 'Справочник статусов';

create table service (
	service_id uuid not null primary key,
	name text not null
);
comment on table service is 'Справочник сервисов';

create table log (
	stamp timestamp without time zone not null,
	channel uuid not null references channel(channel_id),
	region uuid not null references region(region_id),
	service uuid not null references service(service_id),
	status uuid not null references status(status_id),
	"user" bigint not null
);
comment on table log is 'Лог событий из лог-файла';
