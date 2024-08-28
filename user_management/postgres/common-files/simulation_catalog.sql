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
-- Name: simulation_catalog; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.simulation_catalog (
    sim_name character varying NOT NULL,
    build_type character varying,
    sim_url character varying
);


ALTER TABLE public.simulation_catalog OWNER TO postgres;

--
-- Data for Name: simulation_catalog; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.simulation_catalog (sim_name, build_type, sim_url) FROM stdin;
Sim_2K_cells	user	https://arm.seli.gic.ericsson.se/artifactory/proj-eric-oss-restsim-ci-internal-helm-local/eric-oss-restsim-release/eric-oss-restsim-release-1.0.0-56.tgz
Sim_5K_cells	dev	https://arm.seli.gic.ericsson.se/artifactory/proj-eric-oss-restsim-ci-internal-helm-local/eric-oss-restsim-release/eric-oss-restsim-release-1.0.0-57.tgz
Sim_80K_cells	user	https://arm.seli.gic.ericsson.se/artifactory/proj-eric-oss-restsim-ci-internal-helm-local/eric-oss-restsim-release/eric-oss-restsim-release-1.0.0-58.tgz
\.


--
-- Name: simulation_catalog simulation_catalog_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.simulation_catalog
    ADD CONSTRAINT simulation_catalog_pkey PRIMARY KEY (sim_name);


--
-- PostgreSQL database dump complete
--

