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
-- Name: user_management; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_management (
    signum character varying NOT NULL,
    team_name character varying,
    role character varying
);


ALTER TABLE public.user_management OWNER TO postgres;

--
-- Data for Name: user_management; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_management (signum, team_name, role) FROM stdin;
\.


--
-- Name: user_management user_management_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_management
    ADD CONSTRAINT user_management_pkey PRIMARY KEY (signum);


--
-- PostgreSQL database dump complete
--

