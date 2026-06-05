ALTER TABLE health_states
ADD CONSTRAINT chk_health_states_current_lte_max
CHECK (current_hp <= max_hp);
