--
-- PostgreSQL database dump
--

-- Dumped from database version 14.4
-- Dumped by pg_dump version 14.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: auto_text_fields; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.auto_text_fields (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    module_name text NOT NULL,
    item_name text NOT NULL,
    text text NOT NULL,
    organization_id uuid
);


ALTER TABLE public.auto_text_fields OWNER TO liyanxian;

--
-- Name: basic_charge_settings; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.basic_charge_settings (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    sort_index bigint,
    organization_id uuid,
    organization_basic_charge_setting_id uuid,
    patient_id uuid,
    user_id uuid
);


ALTER TABLE public.basic_charge_settings OWNER TO liyanxian;

--
-- Name: basic_charges; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.basic_charges (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    item_name text,
    type text,
    unit text,
    price bigint,
    tax_type text,
    start_date timestamp with time zone,
    end_date timestamp with time zone,
    note text,
    receipt_status text,
    receipt_date timestamp with time zone,
    sort_index bigint,
    organization_id uuid,
    patient_id uuid,
    user_id uuid
);


ALTER TABLE public.basic_charges OWNER TO liyanxian;

--
-- Name: deposit_records; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.deposit_records (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    id_number text NOT NULL,
    date timestamp with time zone,
    type text NOT NULL,
    price bigint NOT NULL,
    drawee text,
    note text,
    invalid boolean NOT NULL,
    user_id uuid,
    patient_id uuid,
    organization_id uuid
);


ALTER TABLE public.deposit_records OWNER TO liyanxian;

--
-- Name: files; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.files (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    file_name text,
    url text,
    organization_id uuid
);


ALTER TABLE public.files OWNER TO liyanxian;

--
-- Name: non_fixed_charge_records; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.non_fixed_charge_records (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    non_fixed_charge_date timestamp with time zone NOT NULL,
    item_category text NOT NULL,
    item_name text NOT NULL,
    type text NOT NULL,
    unit text NOT NULL,
    price bigint NOT NULL,
    quantity bigint NOT NULL,
    subtotal bigint NOT NULL,
    note text,
    is_tax text NOT NULL,
    tax_type text,
    receipt_status text,
    receipt_date timestamp with time zone,
    organization_id uuid,
    patient_id uuid,
    user_id uuid
);


ALTER TABLE public.non_fixed_charge_records OWNER TO liyanxian;

--
-- Name: organization_basic_charge_settings; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.organization_basic_charge_settings (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    item_name text NOT NULL,
    type text NOT NULL,
    unit text NOT NULL,
    price bigint NOT NULL,
    is_tax text NOT NULL,
    tax_type text,
    organization_id uuid
);


ALTER TABLE public.organization_basic_charge_settings OWNER TO liyanxian;

--
-- Name: organization_non_fixed_charge_settings; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.organization_non_fixed_charge_settings (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    item_category text NOT NULL,
    item_name text NOT NULL,
    type text NOT NULL,
    unit text NOT NULL,
    price bigint NOT NULL,
    is_tax text NOT NULL,
    tax_type text,
    organization_id uuid
);


ALTER TABLE public.organization_non_fixed_charge_settings OWNER TO liyanxian;

--
-- Name: organization_receipt_template_settings; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.organization_receipt_template_settings (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    tax_types text[],
    organization_picture text,
    title_name text NOT NULL,
    patient_info text[],
    price_show_type text NOT NULL,
    organization_info_one text[],
    organization_info_two text[],
    note_text text,
    seal_one_name text,
    seal_one_picture text,
    seal_two_name text,
    seal_two_picture text,
    seal_three_name text,
    seal_three_picture text,
    seal_four_name text,
    seal_four_picture text,
    part_one_name text NOT NULL,
    part_two_name text NOT NULL,
    organization_id uuid
);


ALTER TABLE public.organization_receipt_template_settings OWNER TO liyanxian;

--
-- Name: organization_receipts; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.organization_receipts (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    first_text text,
    year text,
    year_text text,
    month text,
    month_text text,
    last_text text,
    is_reset_in_next_cycle boolean,
    organization_id uuid
);


ALTER TABLE public.organization_receipts OWNER TO liyanxian;

--
-- Name: organizations; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.organizations (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    address_city text,
    address_district text,
    address text,
    phone text,
    fax text,
    owner text,
    email text,
    tax_id_number text,
    remittance_bank text,
    remittance_id_number text,
    remittance_user_name text,
    establishment_number text,
    solution text,
    fixed_charge_start_month text NOT NULL,
    fixed_charge_start_date bigint NOT NULL,
    fixed_charge_end_month text NOT NULL,
    fixed_charge_end_date bigint NOT NULL,
    non_fixed_charge_start_month text NOT NULL,
    non_fixed_charge_start_date bigint NOT NULL,
    non_fixed_charge_end_month text NOT NULL,
    non_fixed_charge_end_date bigint NOT NULL,
    transfer_refund_start_month text NOT NULL,
    transfer_refund_start_date bigint NOT NULL,
    transfer_refund_end_month text NOT NULL,
    transfer_refund_end_date bigint NOT NULL,
    branchs text[],
    provider_org_id text,
    privacy text,
    test_time timestamp with time zone
);


ALTER TABLE public.organizations OWNER TO liyanxian;

--
-- Name: patient_bill_basic_charges; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.patient_bill_basic_charges (
    patient_bill_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    basic_charge_id uuid DEFAULT public.uuid_generate_v4() NOT NULL
);


ALTER TABLE public.patient_bill_basic_charges OWNER TO liyanxian;

--
-- Name: patient_bill_non_fixed_charge_records; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.patient_bill_non_fixed_charge_records (
    patient_bill_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    non_fixed_charge_record_id uuid DEFAULT public.uuid_generate_v4() NOT NULL
);


ALTER TABLE public.patient_bill_non_fixed_charge_records OWNER TO liyanxian;

--
-- Name: patient_bill_subsidies; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.patient_bill_subsidies (
    patient_bill_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    subsidy_id uuid DEFAULT public.uuid_generate_v4() NOT NULL
);


ALTER TABLE public.patient_bill_subsidies OWNER TO liyanxian;

--
-- Name: patient_bill_transfer_refund_leaves; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.patient_bill_transfer_refund_leaves (
    patient_bill_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    transfer_refund_leave_id uuid DEFAULT public.uuid_generate_v4() NOT NULL
);


ALTER TABLE public.patient_bill_transfer_refund_leaves OWNER TO liyanxian;

--
-- Name: patient_bills; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.patient_bills (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    amount_received bigint,
    note text,
    edit_note_date timestamp with time zone,
    fixed_charge_start_date timestamp with time zone,
    fixed_charge_end_date timestamp with time zone,
    transfer_refund_start_date timestamp with time zone,
    transfer_refund_end_date timestamp with time zone,
    non_fixed_charge_start_date timestamp with time zone,
    non_fixed_charge_end_date timestamp with time zone,
    basic_charges_sort_ids uuid[],
    subsidies_sort_ids uuid[],
    bill_year bigint,
    bill_month bigint,
    organization_id uuid,
    patient_id uuid,
    user_id uuid,
    edit_note_user_id uuid,
    amount_due bigint
);


ALTER TABLE public.patient_bills OWNER TO liyanxian;

--
-- Name: patient_user_relation; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.patient_user_relation (
    patient_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid DEFAULT public.uuid_generate_v4() NOT NULL
);


ALTER TABLE public.patient_user_relation OWNER TO liyanxian;

--
-- Name: patients; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.patients (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    first_name text,
    last_name text,
    id_number text,
    photo_url text,
    photo_x_position bigint,
    photo_y_position bigint,
    provider_id text,
    status text,
    branch text,
    room text,
    bed text,
    sex text,
    birthday timestamp with time zone,
    check_in_date timestamp with time zone,
    patient_number text,
    record_number text,
    numbering text,
    organization_id uuid
);


ALTER TABLE public.patients OWNER TO liyanxian;

--
-- Name: pay_record_details; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.pay_record_details (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    record_date timestamp with time zone,
    type text NOT NULL,
    price bigint NOT NULL,
    method text NOT NULL,
    payer text,
    handler text,
    note text,
    pay_record_id uuid,
    organization_id uuid,
    patient_id uuid,
    user_id uuid
);


ALTER TABLE public.pay_record_details OWNER TO liyanxian;

--
-- Name: pay_record_pay_record_details; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.pay_record_pay_record_details (
    pay_record_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    pay_record_detail_id uuid DEFAULT public.uuid_generate_v4() NOT NULL
);


ALTER TABLE public.pay_record_pay_record_details OWNER TO liyanxian;

--
-- Name: pay_records; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.pay_records (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    basic_charge jsonb,
    subsidy jsonb,
    transfer_refund_hospitalized jsonb,
    transfer_refund_leave jsonb,
    transfer_refund_discharge jsonb,
    non_fixed_charge jsonb,
    pay_date timestamp with time zone,
    receipt_number text,
    tax_type text,
    amount_due bigint,
    paid_amount bigint,
    note text,
    is_invalid boolean,
    invalid_caption text,
    invalid_date timestamp with time zone,
    pay_year bigint,
    pay_month bigint,
    patient_bill_id uuid,
    organization_id uuid,
    patient_id uuid,
    user_id uuid,
    created_user_id uuid,
    invalid_user_id uuid
);


ALTER TABLE public.pay_records OWNER TO liyanxian;

--
-- Name: subsidies; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.subsidies (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    item_name text NOT NULL,
    type text NOT NULL,
    price bigint NOT NULL,
    unit text,
    id_number text,
    note text,
    start_date timestamp with time zone,
    end_date timestamp with time zone,
    receipt_status text,
    receipt_date timestamp with time zone,
    sort_index bigint,
    organization_id uuid,
    patient_id uuid,
    user_id uuid
);


ALTER TABLE public.subsidies OWNER TO liyanxian;

--
-- Name: subsidy_settings; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.subsidy_settings (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    sort_index bigint,
    item_name text NOT NULL,
    type text NOT NULL,
    price bigint NOT NULL,
    unit text,
    id_number text,
    note text,
    organization_id uuid,
    patient_id uuid
);


ALTER TABLE public.subsidy_settings OWNER TO liyanxian;

--
-- Name: transfer_refund_leaves; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.transfer_refund_leaves (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone NOT NULL,
    reason text NOT NULL,
    is_reserve_bed text NOT NULL,
    note text,
    items jsonb,
    receipt_status text,
    receipt_date timestamp with time zone,
    organization_id uuid,
    patient_id uuid,
    user_id uuid
);


ALTER TABLE public.transfer_refund_leaves OWNER TO liyanxian;

--
-- Name: users; Type: TABLE; Schema: public; Owner: liyanxian
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    first_name text,
    last_name text,
    display_name text,
    id_number text,
    preference jsonb,
    token text,
    token_expired_at timestamp with time zone,
    provider_token jsonb,
    provider_id text,
    username text,
    password text,
    organization_id uuid
);


ALTER TABLE public.users OWNER TO liyanxian;

--
-- Data for Name: auto_text_fields; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.auto_text_fields (id, created_at, updated_at, deleted_at, module_name, item_name, text, organization_id) FROM stdin;
\.


--
-- Data for Name: basic_charge_settings; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.basic_charge_settings (id, created_at, updated_at, deleted_at, sort_index, organization_id, organization_basic_charge_setting_id, patient_id, user_id) FROM stdin;
\.


--
-- Data for Name: basic_charges; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.basic_charges (id, created_at, updated_at, deleted_at, item_name, type, unit, price, tax_type, start_date, end_date, note, receipt_status, receipt_date, sort_index, organization_id, patient_id, user_id) FROM stdin;
\.


--
-- Data for Name: deposit_records; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.deposit_records (id, created_at, updated_at, deleted_at, id_number, date, type, price, drawee, note, invalid, user_id, patient_id, organization_id) FROM stdin;
\.


--
-- Data for Name: files; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.files (id, created_at, updated_at, deleted_at, file_name, url, organization_id) FROM stdin;
\.


--
-- Data for Name: non_fixed_charge_records; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.non_fixed_charge_records (id, created_at, updated_at, deleted_at, non_fixed_charge_date, item_category, item_name, type, unit, price, quantity, subtotal, note, is_tax, tax_type, receipt_status, receipt_date, organization_id, patient_id, user_id) FROM stdin;
\.


--
-- Data for Name: organization_basic_charge_settings; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.organization_basic_charge_settings (id, created_at, updated_at, deleted_at, item_name, type, unit, price, is_tax, tax_type, organization_id) FROM stdin;
\.


--
-- Data for Name: organization_non_fixed_charge_settings; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.organization_non_fixed_charge_settings (id, created_at, updated_at, deleted_at, item_category, item_name, type, unit, price, is_tax, tax_type, organization_id) FROM stdin;
\.


--
-- Data for Name: organization_receipt_template_settings; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.organization_receipt_template_settings (id, created_at, updated_at, deleted_at, name, tax_types, organization_picture, title_name, patient_info, price_show_type, organization_info_one, organization_info_two, note_text, seal_one_name, seal_one_picture, seal_two_name, seal_two_picture, seal_three_name, seal_three_picture, seal_four_name, seal_four_picture, part_one_name, part_two_name, organization_id) FROM stdin;
\.


--
-- Data for Name: organization_receipts; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.organization_receipts (id, created_at, updated_at, deleted_at, first_text, year, year_text, month, month_text, last_text, is_reset_in_next_cycle, organization_id) FROM stdin;
\.


--
-- Data for Name: organizations; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.organizations (id, created_at, updated_at, deleted_at, name, address_city, address_district, address, phone, fax, owner, email, tax_id_number, remittance_bank, remittance_id_number, remittance_user_name, establishment_number, solution, fixed_charge_start_month, fixed_charge_start_date, fixed_charge_end_month, fixed_charge_end_date, non_fixed_charge_start_month, non_fixed_charge_start_date, non_fixed_charge_end_month, non_fixed_charge_end_date, transfer_refund_start_month, transfer_refund_start_date, transfer_refund_end_month, transfer_refund_end_date, branchs, provider_org_id, privacy, test_time) FROM stdin;
a9994b72-c534-423f-b857-615133ffa248	2022-03-16 11:20:55.573546+08	2023-04-07 10:39:51.821315+08	\N	智齡護家	新北市	雙溪區	123號56巷	12345677888	xxxxx	qq1	10	50868012	匯款銀行	匯款帳號	匯款戶名	設立許可文號	nis	thisMonth	1	lastMonth	31	lastMonth	1	lastMonth	31	lastMonth	1	lastMonth	31	{1,2,3,4}	5c10bdf47b43650f407de7d6	true	\N
\.


--
-- Data for Name: patient_bill_basic_charges; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.patient_bill_basic_charges (patient_bill_id, basic_charge_id) FROM stdin;
\.


--
-- Data for Name: patient_bill_non_fixed_charge_records; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.patient_bill_non_fixed_charge_records (patient_bill_id, non_fixed_charge_record_id) FROM stdin;
\.


--
-- Data for Name: patient_bill_subsidies; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.patient_bill_subsidies (patient_bill_id, subsidy_id) FROM stdin;
\.


--
-- Data for Name: patient_bill_transfer_refund_leaves; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.patient_bill_transfer_refund_leaves (patient_bill_id, transfer_refund_leave_id) FROM stdin;
\.


--
-- Data for Name: patient_bills; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.patient_bills (id, created_at, updated_at, deleted_at, amount_received, note, edit_note_date, fixed_charge_start_date, fixed_charge_end_date, transfer_refund_start_date, transfer_refund_end_date, non_fixed_charge_start_date, non_fixed_charge_end_date, basic_charges_sort_ids, subsidies_sort_ids, bill_year, bill_month, organization_id, patient_id, user_id, edit_note_user_id, amount_due) FROM stdin;
\.


--
-- Data for Name: patient_user_relation; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.patient_user_relation (patient_id, user_id) FROM stdin;
\.


--
-- Data for Name: patients; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.patients (id, created_at, updated_at, deleted_at, first_name, last_name, id_number, photo_url, photo_x_position, photo_y_position, provider_id, status, branch, room, bed, sex, birthday, check_in_date, patient_number, record_number, numbering, organization_id) FROM stdin;
317edd46-1b4e-418c-8395-38f4eff59a8c	2023-03-25 14:49:30.061221+08	2023-03-25 14:50:00.865647+08	\N	可可	林	V128457392	https://storage.googleapis.com/public-origin-jubo-image/patient/5c10bdf47b43650f407de7d6/08d9592a3fb0926d7fd1af616f573015	1	1	62c24a58548e2a00285cffa4	present	1	103	B	male	1932-10-15 09:58:00+08	2022-07-08 15:05:21.365+08	\N	\N	\N	a9994b72-c534-423f-b857-615133ffa248
\.


--
-- Data for Name: pay_record_details; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.pay_record_details (id, created_at, updated_at, deleted_at, record_date, type, price, method, payer, handler, note, pay_record_id, organization_id, patient_id, user_id) FROM stdin;
\.


--
-- Data for Name: pay_record_pay_record_details; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.pay_record_pay_record_details (pay_record_id, pay_record_detail_id) FROM stdin;
\.


--
-- Data for Name: pay_records; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.pay_records (id, created_at, updated_at, deleted_at, basic_charge, subsidy, transfer_refund_hospitalized, transfer_refund_leave, transfer_refund_discharge, non_fixed_charge, pay_date, receipt_number, tax_type, amount_due, paid_amount, note, is_invalid, invalid_caption, invalid_date, pay_year, pay_month, patient_bill_id, organization_id, patient_id, user_id, created_user_id, invalid_user_id) FROM stdin;
\.


--
-- Data for Name: subsidies; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.subsidies (id, created_at, updated_at, deleted_at, item_name, type, price, unit, id_number, note, start_date, end_date, receipt_status, receipt_date, sort_index, organization_id, patient_id, user_id) FROM stdin;
\.


--
-- Data for Name: subsidy_settings; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.subsidy_settings (id, created_at, updated_at, deleted_at, sort_index, item_name, type, price, unit, id_number, note, organization_id, patient_id) FROM stdin;
\.


--
-- Data for Name: transfer_refund_leaves; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.transfer_refund_leaves (id, created_at, updated_at, deleted_at, start_date, end_date, reason, is_reserve_bed, note, items, receipt_status, receipt_date, organization_id, patient_id, user_id) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: liyanxian
--

COPY public.users (id, created_at, updated_at, deleted_at, first_name, last_name, display_name, id_number, preference, token, token_expired_at, provider_token, provider_id, username, password, organization_id) FROM stdin;
c71f802b-fea7-4278-9062-1181937df0d9	2022-03-16 11:20:55.742209+08	2023-09-18 15:30:47.260418+08	\N	小芳	長	長小芳	S123333444	{"era": "Christian", "staff": "all", "branch": "all", "locale": "zhTW", "healthDashboardGroups": [{"items": [{"items": [{"name": "customizedFormLine?formId=6476df5e439112d75fcfff51&version=6476df7d439112d75fcfff87&startField=1&valueField=2", "width": 2, "contents": [{"name": "DepressionLineChart", "height": 180}, {"name": "Depression5LineChart", "height": 180}, {"name": "customizedFormFloorLine?formId=6476df5e439112d75fcfff51&version=6476df7d439112d75fcfff87&startField=1&endField=3&valueField=2&preprocessing=oneLine&chartType=line", "height": 180}, {"name": "customizedFormFloorLine?formId=6476df5e439112d75fcfff51&version=6476df7d439112d75fcfff87&startField=1&endField=3&valueField=2", "height": 180}, {"name": "customizedFormTimeBar?formId=647837056992d64eca900257&version=647837b66992d64eca900291&startField=1&endField=2&valueField=3&noteField=4&showOther=false", "height": 180}, {"name": "customizedFormTimeBar?formId=647989eb1c98e60df880f0b4&version=64798a5b1c98e60df880f0ee&startField=1&endField=4&valueField=2&colors=#ffffff,#bdbdbd", "height": 180}]}, {"name": "TPRChart_T", "width": 2, "contents": [{"name": "TPRChart_T", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "TPRChart_R", "width": 1, "contents": [{"name": "TPRChart_R", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "TPRChart_P", "width": 1, "contents": [{"name": "TPRChart_P", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "bloodPressureChartSYS", "width": 1, "contents": [{"name": "bloodPressureChartSYS", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "bloodPressureChartDIA", "width": 1, "contents": [{"name": "bloodPressureChartDIA", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}], "range": [1, "hour"]}, {"items": [{"name": "bloodSugarChart", "width": 2, "contents": [{"name": "bloodSugarChart", "height": 150, "hideXAxisText": "true"}, {"name": "insulinChart", "height": 60}]}, {"name": "ioChartTotal", "width": 2, "contents": [{"name": "ioChartTotal", "height": 150, "hideXAxisText": "true"}, {"name": "ioBidirectionalHistogramChart", "height": 60}]}, {"name": "grandmaGuo", "width": 2, "contents": [{"name": "grandmaGuo", "height": 254}]}, {"name": "examReport2Table?title=檢驗報告", "width": 2, "contents": [{"name": "examReport2Table?valueField=all", "height": 508}]}, {"name": "體重紀錄", "width": 1, "contents": [{"name": "customizedFormLine?formId=61970c8989215900262a957a&version=63183d377dcbbd0041823829&startField=1&valueField=3,4&preprocessing=oneLine&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 150, "hideXAxisText": "true"}, {"name": "customizedFormBidirectionalHistogram?formId=61970c8989215900262a957a&version=63183d377dcbbd0041823829&startField=1&valueField=5", "height": 60}]}, {"name": "腹圍(cm)", "width": 1, "contents": [{"name": "customizedFormLine?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=37&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 272.4}]}], "range": [6, "day"]}, {"items": [{"name": "事件記錄-精神狀態", "width": 2, "contents": [{"name": "customizedFormLine?formId=6232fa3878ec040027ed86d0&version=625d1a40a6860e0027ace1c8&startField=3&valueField=4&chartTypes=line&noteField=5", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=635a2b42564b73004104dd5b&version=635a2bee564b73004104f843&startField=2&endField=3&valueField=4&colors=#ffffff,#e0e0e0&noteField=5,6", "height": 30, "hiddenLegend": "true"}]}, {"name": "HRV-LF/HF ratio", "width": 1, "contents": [{"name": "customizedFormLastOne?formId=624d07fdf1eb12002764edb2&version=62691fcd8b4ea700275f4591&startField=3&valueField=7", "height": 254}]}, {"name": "事件紀錄-每日總計", "width": 1, "contents": [{"name": "customizedFormList?formId=625e714294d6280027e6614a&version=625e8158a6860e0027cd5a40&startField=1&valueField=2,3,4", "height": 254}]}, {"name": "HRV-SDNN", "width": 1, "contents": [{"name": "customizedFormLine?formId=624d07fdf1eb12002764edb2&version=62691fcd8b4ea700275f4591&startField=3&valueField=4,9,10,11&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY&chartTypes=lineAndScatter,line,line,line&dasharray=0,5,5,5&colors=#0097A7,#60bfe5,#c0d243,#a683e2", "height": 254}]}, {"name": "HRV-TP", "width": 1, "contents": [{"name": "customizedFormLine?formId=624d07fdf1eb12002764edb2&version=62691fcd8b4ea700275f4591&startField=3&valueField=8&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 254}]}], "range": [2, "day"]}], "label": "全部_醫"}, {"items": [{"items": [{"name": "TPRChart_T", "width": 2, "contents": [{"name": "TPRChart_T", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "TPRChart_R", "width": 1, "contents": [{"name": "TPRChart_R", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "TPRChart_P", "width": 1, "contents": [{"name": "TPRChart_P", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "體重紀錄", "width": 1, "contents": [{"name": "customizedFormLine?formId=61970c8989215900262a957a&version=63183d377dcbbd0041823829&startField=1&valueField=3,4&preprocessing=oneLine&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 150, "hideXAxisText": "true"}, {"name": "customizedFormBidirectionalHistogram?formId=61970c8989215900262a957a&version=63183d377dcbbd0041823829&startField=1&valueField=5", "height": 60}]}, {"name": "examReport2LineChart?title=血檢-血色素 Hemo", "width": 1, "contents": [{"name": "examReport2LineChart?checkUnit=true&valueField=hb&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 254}]}, {"name": "examReport2List?title=血檢-血球容積 Hct", "width": 1, "contents": [{"name": "examReport2List?valueField=ht", "height": 254}]}, {"name": "examReport2List?title=血檢-血小板 Platelet", "width": 1, "contents": [{"name": "examReport2List?valueField=platelets", "height": 254}]}, {"name": "examReport2LineChart?title=生化-血液中尿素氮 BUN", "width": 2, "contents": [{"name": "examReport2LineChart?checkUnit=true&valueField=bun&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 254}]}, {"name": "examReport2LineChart?title=生化-血清肌酸酐 Creatinine", "width": 1, "contents": [{"name": "examReport2LineChart?checkUnit=true&valueField=creatinine", "height": 254}]}, {"name": "examReport2List?title=生化-鈉 NaB", "width": 1, "contents": [{"name": "examReport2List?valueField=sodium", "height": 254}]}, {"name": "examReport2List?title=生化-鉀 KB", "width": 1, "contents": [{"name": "examReport2List?valueField=kalium", "height": 254}]}, {"name": "examReport2List?title=生化-鈣 CaB", "width": 1, "contents": [{"name": "examReport2List?valueField=calcium", "height": 254}]}], "range": [3, "day"]}], "label": "洗腎_醫"}, {"items": [{"items": [{"name": "TPRChart_T", "width": 2, "contents": [{"name": "TPRChart_T", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "TPRChart_R", "width": 1, "contents": [{"name": "TPRChart_R", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "TPRChart_P", "width": 1, "contents": [{"name": "TPRChart_P", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "bloodPressureChart", "width": 1, "contents": [{"name": "bloodPressureChartSYS", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "bloodPressureChart", "width": 1, "contents": [{"name": "bloodPressureChartDIA", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "examReport2LineChart?title=血檢-血色素 Hemo", "width": 1, "contents": [{"name": "examReport2LineChart?checkUnit=true&valueField=hb&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 254}]}, {"name": "examReport2List?title=血檢-白血球 WBCB", "width": 1, "contents": [{"name": "examReport2List?valueField=wbc", "height": 254}]}, {"name": "examReport2List?title=血檢-血小板 Platelet", "width": 1, "contents": [{"name": "examReport2List?valueField=platelets", "height": 254}]}, {"name": "examReport2List?title=血檢-中性球 Neutrophils", "width": 1, "contents": [{"name": "examReport2List?valueField=neutrophils", "height": 254}]}, {"name": "examReport2List?title=血檢-淋巴球 Lympho", "width": 1, "contents": [{"name": "examReport2List?valueField=lymphocytes", "height": 254}]}, {"name": "examReport2List?title=生化-天門冬胺酸轉胺酶 AST_GOT", "width": 1, "contents": [{"name": "examReport2List?valueField=sgot", "height": 254}]}, {"name": "examReport2List?title=生化-丙胺酸轉胺酶 ALT_GPT", "width": 1, "contents": [{"name": "examReport2List?valueField=sgpt", "height": 254}]}], "range": [3, "day"]}], "label": "感染_醫"}, {"items": [{"items": [{"name": "TPRChart_T", "width": 2, "contents": [{"name": "TPRChart_T", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "TPRChart_R", "width": 1, "contents": [{"name": "TPRChart_R", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "TPRChart_P", "width": 1, "contents": [{"name": "TPRChart_P", "height": 180, "hideXAxisText": "true"}, {"name": "customizedFormTimeBar?formId=6233116878ec040027f02231&version=625d1a8aa6860e0027acec82&startField=3&endField=4&valueField=5&noteField=6", "height": 30, "hiddenLegend": "true"}]}, {"name": "體重紀錄", "width": 1, "contents": [{"name": "customizedFormLine?formId=61970c8989215900262a957a&version=63183d377dcbbd0041823829&startField=1&valueField=3,4&preprocessing=oneLine&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 150, "hideXAxisText": "true"}, {"name": "customizedFormBidirectionalHistogram?formId=61970c8989215900262a957a&version=63183d377dcbbd0041823829&startField=1&valueField=5", "height": 60}]}, {"name": "腹圍(cm)", "width": 1, "contents": [{"name": "customizedFormLine?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=37&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 272.4}]}, {"name": "ioChartTotal", "width": 1, "contents": [{"name": "ioChartTotal", "height": 150, "hideXAxisText": "true"}, {"name": "ioBidirectionalHistogramChart", "height": 60}]}, {"name": "examReport2List?title血檢-血球容積 Hct", "width": 1, "contents": [{"name": "examReport2List?checkUnit=true&valueField=ht", "height": 254}]}, {"name": "examReport2LineChart?title=血檢-血色素 Hemo", "width": 2, "contents": [{"name": "examReport2LineChart?checkUnit=true&valueField=hb=&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 254}]}, {"name": "examReport2List?title=血檢-血小板 Platelet", "width": 1, "contents": [{"name": "examReport2List?valueField=platelets", "height": 254}]}, {"name": "examReport2List?title=生化-天門冬胺酸轉胺酶 AST_GOT", "width": 1, "contents": [{"name": "examReport2List?valueField=sgot", "height": 254}]}, {"name": "examReport2List?title=生化-丙胺酸轉胺酶 ALT_GPT", "width": 1, "contents": [{"name": "examReport2List?valueField=sgpt", "height": 254}]}, {"name": "examReport2List?title=生化-血清白蛋白 Albumin", "width": 1, "contents": [{"name": "examReport2List?valueField=albumin", "height": 254}]}, {"name": "examReport2List?title=生化-納 NaB", "width": 1, "contents": [{"name": "examReport2List?valueField=sodium", "height": 254}]}], "range": [3, "day"]}], "label": "腸胃_醫"}, {"items": [{"items": [{"name": "bloodSugarChart", "width": 2, "contents": [{"name": "bloodSugarChart", "height": 150, "hideXAxisText": "true"}, {"name": "insulinChart", "height": 60}]}, {"name": "ioChartTotal", "width": 1, "contents": [{"name": "ioChartTotal", "height": 150, "hideXAxisText": "true"}, {"name": "ioBidirectionalHistogramChart", "height": 60}]}, {"name": "examReport2List?title=生化-鈉 NaB", "width": 1, "contents": [{"name": "examReport2List?valueField=sodium", "height": 254}]}], "range": [3, "day"]}], "label": "血糖_醫"}, {"items": [{"items": [{"name": "TPRChart_T", "width": 2, "contents": [{"name": "TPRChart_T", "height": 180}]}, {"name": "TPRChart_R", "width": 1, "contents": [{"name": "TPRChart_R", "height": 180}]}, {"name": "TPRChart_P", "width": 1, "contents": [{"name": "TPRChart_P", "height": 180}]}, {"name": "bloodPressureChart", "width": 1, "contents": [{"name": "bloodPressureChartSYS", "height": 180}]}, {"name": "bloodPressureChart", "width": 1, "contents": [{"name": "bloodPressureChartDIA", "height": 180}]}, {"name": "ioChartTotal", "width": 2, "contents": [{"name": "ioChartTotal", "height": 150, "hideXAxisText": "true"}, {"name": "ioBidirectionalHistogramChart", "height": 60}]}, {"name": "bloodSugarChart", "width": 2, "contents": [{"name": "bloodSugarChart", "height": 150, "hideXAxisText": "true"}, {"name": "insulinChart", "height": 60}]}], "range": [2, "day"]}, {"items": [{"name": "體重紀錄", "width": 1, "contents": [{"name": "customizedFormLine?formId=61970c8989215900262a957a&version=63183d377dcbbd0041823829&startField=1&valueField=3,4&preprocessing=oneLine&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 150, "hideXAxisText": "true"}, {"name": "customizedFormBidirectionalHistogram?formId=61970c8989215900262a957a&version=63183d377dcbbd0041823829&startField=1&valueField=5", "height": 60}]}, {"name": "腹圍(cm)", "width": 1, "contents": [{"name": "customizedFormLine?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=37&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 272.4}]}, {"name": "examReport2LineChart?title=生化-血液中尿素氮 BUN", "width": 1, "contents": [{"name": "examReport2LineChart?checkUnit=true&valueField=bun&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 254}]}, {"name": "examReport2LineChart?title=生化-鉀 KB", "width": 1, "contents": [{"name": "examReport2LineChart?valueField=kalium", "height": 254}]}, {"name": "examReport2LineChart?title=生化-磷 PB", "width": 2, "contents": [{"name": "examReport2LineChart?valueField=serumPhosphorus", "height": 254}]}, {"name": "examReport2Table?title=血液生化數值", "width": 2, "contents": [{"name": "examReport2Table?valueField=all", "height": 254}]}, {"name": "營養配方", "width": 2, "contents": [{"name": "customizedFormDetailList?formId=6188b45833ff890027ec7b30&version=62a2f7421df9b50027f41de5&startField=1&valueField=2,3,4,5,6,7,8,9,10&column=2", "height": 254}]}, {"name": "grandmaGuo", "width": 2, "contents": [{"name": "grandmaGuo", "height": 254}]}], "range": [6, "day"]}], "label": "全部_營"}, {"items": [{"items": [{"name": "心", "width": 2, "contents": [{"name": "customizedFormDetailList?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=6,7,8,9,10,11,12", "height": 240}]}, {"name": "肺 [呼吸型態]", "width": 2, "contents": [{"name": "customizedFormDetailList?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=14,15,16,17,18,19,20,21,22,23", "height": 240}]}, {"name": "腦[面安指數]", "width": 2, "contents": [{"name": "customizedFormDetailList?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=25,26,27,28,29", "height": 240}]}, {"name": "腎", "width": 2, "contents": [{"name": "customizedFormDetailList?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=31,32,33&column=1", "height": 240}]}, {"name": "腸", "width": 2, "contents": [{"name": "customizedFormDetailList?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=35&column=1", "height": 240}]}, {"name": "[腹部評估]", "width": 2, "contents": [{"name": "customizedFormDetailList?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=37,38,39,40,41", "height": 240}]}, {"name": "樂 [情緒型態]", "width": 2, "contents": [{"name": "customizedFormDetailList?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=43,44,45,46,47&column=1", "height": 240}]}, {"name": "養 [進食情形]", "width": 2, "contents": [{"name": "customizedFormDetailList?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=49,50,51,52&column=1", "height": 240}]}, {"name": "健", "width": 2, "contents": [{"name": "customizedFormDetailList?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=54&column=1", "height": 240}]}, {"name": "眠 [睡眠]", "width": 2, "contents": [{"name": "customizedFormDetailList?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=56,57,58,59,60,61,62", "height": 240}]}, {"name": "grandmaGuo", "width": 2, "contents": [{"name": "grandmaGuo", "height": 254}]}], "range": [1, "day"]}], "label": "交班_護"}, {"items": [{"items": [{"name": "TPRChart_T", "width": 2, "contents": [{"name": "TPRChart_T", "height": 180}]}, {"name": "TPRChart_R", "width": 1, "contents": [{"name": "TPRChart_R", "height": 180}]}, {"name": "TPRChart_P", "width": 1, "contents": [{"name": "TPRChart_P", "height": 180}]}, {"name": "bloodPressureChart", "width": 1, "contents": [{"name": "bloodPressureChartSYS", "height": 180}]}, {"name": "bloodPressureChart", "width": 1, "contents": [{"name": "bloodPressureChartDIA", "height": 180}]}], "range": [0, "day"]}, {"items": [{"name": "bloodSugarChart", "width": 2, "contents": [{"name": "bloodSugarChart", "height": 150, "hideXAxisText": "true"}, {"name": "insulinChart", "height": 60}]}], "range": [2, "day"]}, {"items": [{"name": "體重紀錄", "width": 1, "contents": [{"name": "customizedFormLine?formId=61970c8989215900262a957a&version=63183d377dcbbd0041823829&startField=1&valueField=3,4&preprocessing=oneLine&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 150, "hideXAxisText": "true"}, {"name": "customizedFormBidirectionalHistogram?formId=61970c8989215900262a957a&version=63183d377dcbbd0041823829&startField=1&valueField=5", "height": 60}]}, {"name": "腹圍(cm)", "width": 1, "contents": [{"name": "customizedFormLine?formId=5ff3f8cf748bb10027499730&version=61eb670a099c17002893761b&startField=1&valueField=37&yMin=POSITIVE_INFINITY&yMax=NEGATIVE_INFINITY", "height": 272.4}]}, {"name": "ioChartTotal", "width": 2, "contents": [{"name": "ioChartTotal", "height": 150, "hideXAxisText": "true"}, {"name": "ioBidirectionalHistogramChart", "height": 60}]}, {"name": "ioOutputTable", "width": 2, "contents": [{"name": "ioOutputTable", "height": 254}]}, {"name": "examReport2Table?title=檢驗報告", "width": 2, "contents": [{"name": "examReport2Table?valueField=all", "height": 254}]}, {"name": "grandmaGuo", "width": 2, "contents": [{"name": "grandmaGuo", "height": 254}]}], "range": [6, "day"]}], "label": "全部_護"}]}	2023-09-18 15:30:47.246996+08	2023-09-18 15:30:47.246996+08	{"expiry": "2023-09-18T14:24:23.331422+08:00", "token_type": "JWT", "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjVjMTM0ZGY4ZDkwMGQ3MDAyMTUzZDNiMyIsInJvbGVzIjpbIm9yZ2FuaXphdGlvbi1tYW5hZ2VyIiwiZGlyZWN0b3IiLCJkaWV0aXRpYW4iLCJzeXN0ZW0tYWRtaW4iLCJkZWFuIiwibnVyc2UtcHJhY3RpdGlvbmVyIl0sInNjb3BlIjpbImFsbCJdLCJpYXQiOjE2OTUwMTgxNDMsImV4cCI6MTY5NTA2MTM0M30.4OGaLHOj4I2Mi5hM2EuhYJHwp_Kgn-XG_fHzlyk7xVw", "refresh_token": "nx8sfckpmq5clavnpvm9"}	5c134df8d900d7002153d3b3	\N	a9994b72-c534-423f-b857-615133ffa248	a9994b72-c534-423f-b857-615133ffa248
\.


--
-- Name: auto_text_fields auto_text_fields_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.auto_text_fields
    ADD CONSTRAINT auto_text_fields_pkey PRIMARY KEY (id);


--
-- Name: basic_charge_settings basic_charge_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.basic_charge_settings
    ADD CONSTRAINT basic_charge_settings_pkey PRIMARY KEY (id);


--
-- Name: basic_charges basic_charges_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.basic_charges
    ADD CONSTRAINT basic_charges_pkey PRIMARY KEY (id);


--
-- Name: deposit_records deposit_records_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.deposit_records
    ADD CONSTRAINT deposit_records_pkey PRIMARY KEY (id);


--
-- Name: files files_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.files
    ADD CONSTRAINT files_pkey PRIMARY KEY (id);


--
-- Name: non_fixed_charge_records non_fixed_charge_records_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.non_fixed_charge_records
    ADD CONSTRAINT non_fixed_charge_records_pkey PRIMARY KEY (id);


--
-- Name: organization_basic_charge_settings organization_basic_charge_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.organization_basic_charge_settings
    ADD CONSTRAINT organization_basic_charge_settings_pkey PRIMARY KEY (id);


--
-- Name: organization_non_fixed_charge_settings organization_non_fixed_charge_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.organization_non_fixed_charge_settings
    ADD CONSTRAINT organization_non_fixed_charge_settings_pkey PRIMARY KEY (id);


--
-- Name: organization_receipt_template_settings organization_receipt_template_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.organization_receipt_template_settings
    ADD CONSTRAINT organization_receipt_template_settings_pkey PRIMARY KEY (id);


--
-- Name: organization_receipts organization_receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.organization_receipts
    ADD CONSTRAINT organization_receipts_pkey PRIMARY KEY (id);


--
-- Name: organizations organizations_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_pkey PRIMARY KEY (id);


--
-- Name: patient_bill_basic_charges patient_bill_basic_charges_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_basic_charges
    ADD CONSTRAINT patient_bill_basic_charges_pkey PRIMARY KEY (patient_bill_id, basic_charge_id);


--
-- Name: patient_bill_non_fixed_charge_records patient_bill_non_fixed_charge_records_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_non_fixed_charge_records
    ADD CONSTRAINT patient_bill_non_fixed_charge_records_pkey PRIMARY KEY (patient_bill_id, non_fixed_charge_record_id);


--
-- Name: patient_bill_subsidies patient_bill_subsidies_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_subsidies
    ADD CONSTRAINT patient_bill_subsidies_pkey PRIMARY KEY (patient_bill_id, subsidy_id);


--
-- Name: patient_bill_transfer_refund_leaves patient_bill_transfer_refund_leaves_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_transfer_refund_leaves
    ADD CONSTRAINT patient_bill_transfer_refund_leaves_pkey PRIMARY KEY (patient_bill_id, transfer_refund_leave_id);


--
-- Name: patient_bills patient_bills_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bills
    ADD CONSTRAINT patient_bills_pkey PRIMARY KEY (id);


--
-- Name: patient_user_relation patient_user_relation_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_user_relation
    ADD CONSTRAINT patient_user_relation_pkey PRIMARY KEY (patient_id, user_id);


--
-- Name: patients patients_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patients
    ADD CONSTRAINT patients_pkey PRIMARY KEY (id);


--
-- Name: pay_record_details pay_record_details_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_record_details
    ADD CONSTRAINT pay_record_details_pkey PRIMARY KEY (id);


--
-- Name: pay_record_pay_record_details pay_record_pay_record_details_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_record_pay_record_details
    ADD CONSTRAINT pay_record_pay_record_details_pkey PRIMARY KEY (pay_record_id, pay_record_detail_id);


--
-- Name: pay_records pay_records_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_records
    ADD CONSTRAINT pay_records_pkey PRIMARY KEY (id);


--
-- Name: subsidies subsidies_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.subsidies
    ADD CONSTRAINT subsidies_pkey PRIMARY KEY (id);


--
-- Name: subsidy_settings subsidy_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.subsidy_settings
    ADD CONSTRAINT subsidy_settings_pkey PRIMARY KEY (id);


--
-- Name: transfer_refund_leaves transfer_refund_leaves_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.transfer_refund_leaves
    ADD CONSTRAINT transfer_refund_leaves_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_auto_text_fields_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_auto_text_fields_id ON public.auto_text_fields USING btree (id);


--
-- Name: idx_basic_charge_settings_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_basic_charge_settings_id ON public.basic_charge_settings USING btree (id);


--
-- Name: idx_basic_charges_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_basic_charges_id ON public.basic_charges USING btree (id);


--
-- Name: idx_deposit_records_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_deposit_records_id ON public.deposit_records USING btree (id);


--
-- Name: idx_files_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_files_id ON public.files USING btree (id);


--
-- Name: idx_non_fixed_charge_records_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_non_fixed_charge_records_id ON public.non_fixed_charge_records USING btree (id);


--
-- Name: idx_organization_basic_charge_settings_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_organization_basic_charge_settings_id ON public.organization_basic_charge_settings USING btree (id);


--
-- Name: idx_organization_non_fixed_charge_settings_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_organization_non_fixed_charge_settings_id ON public.organization_non_fixed_charge_settings USING btree (id);


--
-- Name: idx_organization_receipt_template_settings_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_organization_receipt_template_settings_id ON public.organization_receipt_template_settings USING btree (id);


--
-- Name: idx_organization_receipts_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_organization_receipts_id ON public.organization_receipts USING btree (id);


--
-- Name: idx_organizations_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_organizations_id ON public.organizations USING btree (id);


--
-- Name: idx_patient_bills_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_patient_bills_id ON public.patient_bills USING btree (id);


--
-- Name: idx_patients_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_patients_id ON public.patients USING btree (id);


--
-- Name: idx_patients_provider_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_patients_provider_id ON public.patients USING btree (provider_id);


--
-- Name: idx_pay_record_details_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_pay_record_details_id ON public.pay_record_details USING btree (id);


--
-- Name: idx_pay_records_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_pay_records_id ON public.pay_records USING btree (id);


--
-- Name: idx_subsidies_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_subsidies_id ON public.subsidies USING btree (id);


--
-- Name: idx_subsidy_settings_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_subsidy_settings_id ON public.subsidy_settings USING btree (id);


--
-- Name: idx_transfer_refund_leaves_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_transfer_refund_leaves_id ON public.transfer_refund_leaves USING btree (id);


--
-- Name: idx_users_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_users_id ON public.users USING btree (id);


--
-- Name: idx_users_provider_id; Type: INDEX; Schema: public; Owner: liyanxian
--

CREATE UNIQUE INDEX idx_users_provider_id ON public.users USING btree (provider_id);


--
-- Name: auto_text_fields fk_auto_text_fields_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.auto_text_fields
    ADD CONSTRAINT fk_auto_text_fields_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: basic_charge_settings fk_basic_charge_settings_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.basic_charge_settings
    ADD CONSTRAINT fk_basic_charge_settings_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: basic_charge_settings fk_basic_charge_settings_organization_basic_charge_setting; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.basic_charge_settings
    ADD CONSTRAINT fk_basic_charge_settings_organization_basic_charge_setting FOREIGN KEY (organization_basic_charge_setting_id) REFERENCES public.organization_basic_charge_settings(id);


--
-- Name: basic_charge_settings fk_basic_charge_settings_patient; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.basic_charge_settings
    ADD CONSTRAINT fk_basic_charge_settings_patient FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: basic_charge_settings fk_basic_charge_settings_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.basic_charge_settings
    ADD CONSTRAINT fk_basic_charge_settings_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: basic_charges fk_basic_charges_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.basic_charges
    ADD CONSTRAINT fk_basic_charges_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: basic_charges fk_basic_charges_patient; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.basic_charges
    ADD CONSTRAINT fk_basic_charges_patient FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: basic_charges fk_basic_charges_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.basic_charges
    ADD CONSTRAINT fk_basic_charges_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: deposit_records fk_deposit_records_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.deposit_records
    ADD CONSTRAINT fk_deposit_records_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: deposit_records fk_deposit_records_patient; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.deposit_records
    ADD CONSTRAINT fk_deposit_records_patient FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: deposit_records fk_deposit_records_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.deposit_records
    ADD CONSTRAINT fk_deposit_records_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: files fk_files_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.files
    ADD CONSTRAINT fk_files_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: non_fixed_charge_records fk_non_fixed_charge_records_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.non_fixed_charge_records
    ADD CONSTRAINT fk_non_fixed_charge_records_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: non_fixed_charge_records fk_non_fixed_charge_records_patient; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.non_fixed_charge_records
    ADD CONSTRAINT fk_non_fixed_charge_records_patient FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: non_fixed_charge_records fk_non_fixed_charge_records_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.non_fixed_charge_records
    ADD CONSTRAINT fk_non_fixed_charge_records_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: organization_basic_charge_settings fk_organization_basic_charge_settings_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.organization_basic_charge_settings
    ADD CONSTRAINT fk_organization_basic_charge_settings_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: organization_non_fixed_charge_settings fk_organization_non_fixed_charge_settings_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.organization_non_fixed_charge_settings
    ADD CONSTRAINT fk_organization_non_fixed_charge_settings_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: organization_receipt_template_settings fk_organization_receipt_template_settings_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.organization_receipt_template_settings
    ADD CONSTRAINT fk_organization_receipt_template_settings_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: organization_receipts fk_organization_receipts_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.organization_receipts
    ADD CONSTRAINT fk_organization_receipts_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: patient_bill_basic_charges fk_patient_bill_basic_charges_basic_charge; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_basic_charges
    ADD CONSTRAINT fk_patient_bill_basic_charges_basic_charge FOREIGN KEY (basic_charge_id) REFERENCES public.basic_charges(id);


--
-- Name: patient_bill_basic_charges fk_patient_bill_basic_charges_patient_bill; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_basic_charges
    ADD CONSTRAINT fk_patient_bill_basic_charges_patient_bill FOREIGN KEY (patient_bill_id) REFERENCES public.patient_bills(id);


--
-- Name: patient_bill_non_fixed_charge_records fk_patient_bill_non_fixed_charge_records_non_fixed_charge_recor; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_non_fixed_charge_records
    ADD CONSTRAINT fk_patient_bill_non_fixed_charge_records_non_fixed_charge_recor FOREIGN KEY (non_fixed_charge_record_id) REFERENCES public.non_fixed_charge_records(id);


--
-- Name: patient_bill_non_fixed_charge_records fk_patient_bill_non_fixed_charge_records_patient_bill; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_non_fixed_charge_records
    ADD CONSTRAINT fk_patient_bill_non_fixed_charge_records_patient_bill FOREIGN KEY (patient_bill_id) REFERENCES public.patient_bills(id);


--
-- Name: patient_bill_subsidies fk_patient_bill_subsidies_patient_bill; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_subsidies
    ADD CONSTRAINT fk_patient_bill_subsidies_patient_bill FOREIGN KEY (patient_bill_id) REFERENCES public.patient_bills(id);


--
-- Name: patient_bill_subsidies fk_patient_bill_subsidies_subsidy; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_subsidies
    ADD CONSTRAINT fk_patient_bill_subsidies_subsidy FOREIGN KEY (subsidy_id) REFERENCES public.subsidies(id);


--
-- Name: patient_bill_transfer_refund_leaves fk_patient_bill_transfer_refund_leaves_patient_bill; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_transfer_refund_leaves
    ADD CONSTRAINT fk_patient_bill_transfer_refund_leaves_patient_bill FOREIGN KEY (patient_bill_id) REFERENCES public.patient_bills(id);


--
-- Name: patient_bill_transfer_refund_leaves fk_patient_bill_transfer_refund_leaves_transfer_refund_leave; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bill_transfer_refund_leaves
    ADD CONSTRAINT fk_patient_bill_transfer_refund_leaves_transfer_refund_leave FOREIGN KEY (transfer_refund_leave_id) REFERENCES public.transfer_refund_leaves(id);


--
-- Name: patient_bills fk_patient_bills_edit_note_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bills
    ADD CONSTRAINT fk_patient_bills_edit_note_user FOREIGN KEY (edit_note_user_id) REFERENCES public.users(id);


--
-- Name: patient_bills fk_patient_bills_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bills
    ADD CONSTRAINT fk_patient_bills_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: patient_bills fk_patient_bills_patient; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bills
    ADD CONSTRAINT fk_patient_bills_patient FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: patient_bills fk_patient_bills_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_bills
    ADD CONSTRAINT fk_patient_bills_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: patient_user_relation fk_patient_user_relation_patient; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_user_relation
    ADD CONSTRAINT fk_patient_user_relation_patient FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: patient_user_relation fk_patient_user_relation_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patient_user_relation
    ADD CONSTRAINT fk_patient_user_relation_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: patients fk_patients_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.patients
    ADD CONSTRAINT fk_patients_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: pay_record_details fk_pay_record_details_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_record_details
    ADD CONSTRAINT fk_pay_record_details_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: pay_record_details fk_pay_record_details_patient; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_record_details
    ADD CONSTRAINT fk_pay_record_details_patient FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: pay_record_details fk_pay_record_details_pay_record; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_record_details
    ADD CONSTRAINT fk_pay_record_details_pay_record FOREIGN KEY (pay_record_id) REFERENCES public.pay_records(id);


--
-- Name: pay_record_details fk_pay_record_details_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_record_details
    ADD CONSTRAINT fk_pay_record_details_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: pay_record_pay_record_details fk_pay_record_pay_record_details_pay_record; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_record_pay_record_details
    ADD CONSTRAINT fk_pay_record_pay_record_details_pay_record FOREIGN KEY (pay_record_id) REFERENCES public.pay_records(id);


--
-- Name: pay_record_pay_record_details fk_pay_record_pay_record_details_pay_record_detail; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_record_pay_record_details
    ADD CONSTRAINT fk_pay_record_pay_record_details_pay_record_detail FOREIGN KEY (pay_record_detail_id) REFERENCES public.pay_record_details(id);


--
-- Name: pay_records fk_pay_records_created_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_records
    ADD CONSTRAINT fk_pay_records_created_user FOREIGN KEY (created_user_id) REFERENCES public.users(id);


--
-- Name: pay_records fk_pay_records_invalid_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_records
    ADD CONSTRAINT fk_pay_records_invalid_user FOREIGN KEY (invalid_user_id) REFERENCES public.users(id);


--
-- Name: pay_records fk_pay_records_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_records
    ADD CONSTRAINT fk_pay_records_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: pay_records fk_pay_records_patient; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_records
    ADD CONSTRAINT fk_pay_records_patient FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: pay_records fk_pay_records_patient_bill; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_records
    ADD CONSTRAINT fk_pay_records_patient_bill FOREIGN KEY (patient_bill_id) REFERENCES public.patient_bills(id);


--
-- Name: pay_records fk_pay_records_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.pay_records
    ADD CONSTRAINT fk_pay_records_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: subsidies fk_subsidies_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.subsidies
    ADD CONSTRAINT fk_subsidies_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: subsidies fk_subsidies_patient; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.subsidies
    ADD CONSTRAINT fk_subsidies_patient FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: subsidies fk_subsidies_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.subsidies
    ADD CONSTRAINT fk_subsidies_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: subsidy_settings fk_subsidy_settings_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.subsidy_settings
    ADD CONSTRAINT fk_subsidy_settings_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: subsidy_settings fk_subsidy_settings_patient; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.subsidy_settings
    ADD CONSTRAINT fk_subsidy_settings_patient FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: transfer_refund_leaves fk_transfer_refund_leaves_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.transfer_refund_leaves
    ADD CONSTRAINT fk_transfer_refund_leaves_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: transfer_refund_leaves fk_transfer_refund_leaves_patient; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.transfer_refund_leaves
    ADD CONSTRAINT fk_transfer_refund_leaves_patient FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: transfer_refund_leaves fk_transfer_refund_leaves_user; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.transfer_refund_leaves
    ADD CONSTRAINT fk_transfer_refund_leaves_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: users fk_users_organization; Type: FK CONSTRAINT; Schema: public; Owner: liyanxian
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_organization FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- PostgreSQL database dump complete
--

