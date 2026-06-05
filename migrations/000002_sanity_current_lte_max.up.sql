ALTER TABLE sanity_states
ADD CONSTRAINT chk_sanity_states_current_lte_max
CHECK (current_sanity <= max_sanity);
