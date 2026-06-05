ALTER TABLE luck_states
ADD CONSTRAINT chk_luck_states_current_lte_starting
CHECK (current_luck <= starting_luck);
