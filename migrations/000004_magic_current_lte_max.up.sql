ALTER TABLE magic_states
ADD CONSTRAINT chk_magic_states_current_lte_max
CHECK (current_mp <= max_mp);
