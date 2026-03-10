-- Remove non-sovereign entities and their associated answer records.
-- Keeps all sovereign nations, plus Taiwan, Palestine, and Kosovo.

DELETE FROM answer
WHERE selected_country IN (SELECT id FROM country WHERE iso2 IN (
    'eu',
    'gb-eng', 'gb-nir', 'gb-sct', 'gb-wls',
    'hk', 'mo',
    'ai', 'as', 'aw', 'ax', 'bm', 'bq', 'cc', 'ck', 'cw', 'cx',
    'fk', 'fo', 'gf', 'gg', 'gi', 'gl', 'gp', 'gu',
    'im', 'io', 'je', 'ky',
    'mp', 'ms', 'nf', 'nu',
    'pf', 'pm', 'pn', 'pr',
    're', 'sh', 'sj', 'sx',
    'tc', 'tk', 'vg', 'vi', 'wf',
    'aq', 'bv', 'gs', 'hm', 'tf',
    'eh'
))
OR correct_country IN (SELECT id FROM country WHERE iso2 IN (
    'eu',
    'gb-eng', 'gb-nir', 'gb-sct', 'gb-wls',
    'hk', 'mo',
    'ai', 'as', 'aw', 'ax', 'bm', 'bq', 'cc', 'ck', 'cw', 'cx',
    'fk', 'fo', 'gf', 'gg', 'gi', 'gl', 'gp', 'gu',
    'im', 'io', 'je', 'ky',
    'mp', 'ms', 'nf', 'nu',
    'pf', 'pm', 'pn', 'pr',
    're', 'sh', 'sj', 'sx',
    'tc', 'tk', 'vg', 'vi', 'wf',
    'aq', 'bv', 'gs', 'hm', 'tf',
    'eh'
));

DELETE FROM country WHERE iso2 IN (
    'eu',
    'gb-eng', 'gb-nir', 'gb-sct', 'gb-wls',
    'hk', 'mo',
    'ai', 'as', 'aw', 'ax', 'bm', 'bq', 'cc', 'ck', 'cw', 'cx',
    'fk', 'fo', 'gf', 'gg', 'gi', 'gl', 'gp', 'gu',
    'im', 'io', 'je', 'ky',
    'mp', 'ms', 'nf', 'nu',
    'pf', 'pm', 'pn', 'pr',
    're', 'sh', 'sj', 'sx',
    'tc', 'tk', 'vg', 'vi', 'wf',
    'aq', 'bv', 'gs', 'hm', 'tf',
    'eh'
);
