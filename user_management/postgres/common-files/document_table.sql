--
-- PostgreSQL database dump
--

-- Dumped from database version 11.18
-- Dumped by pg_dump version 11.18

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

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: document_table; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.document_table (
    document_name character varying NOT NULL,
    document_link character varying
);


ALTER TABLE public.document_table OWNER TO postgres;

--
-- Data for Name: document_table; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.document_table (document_name, document_link) FROM stdin;
Product Documentation	https://arm1s11-eiffel004.eiffel.gic.ericsson.se:8443/nexus/content/sites/tor/idun-sdk/latest/index.html#home
User Onboarding Hyperlink	https://arm1s11-eiffel004.eiffel.gic.ericsson.se:8443/nexus/content/sites/tor/idun-sdk/latest/index.html#home
\.


--
-- Name: document_table document_table_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.document_table
    ADD CONSTRAINT document_table_pkey PRIMARY KEY (document_name);


--
-- PostgreSQL database dump complete
--

