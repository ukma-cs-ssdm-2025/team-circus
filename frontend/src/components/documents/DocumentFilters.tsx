import { FormControl, InputLabel, MenuItem, Select, Stack, TextField } from '@mui/material';
import type { GroupOption, DocumentFilters as DocumentFiltersState } from '../../types/entities';

interface DocumentFiltersProps {
  filters: DocumentFiltersState;
  groupOptions: GroupOption[];
  onGroupChange: (value: string) => void;
  onSearchChange: (value: string) => void;
  filterGroupLabel: string;
  filterAllLabel: string;
  searchPlaceholder: string;
}

const DocumentFilters = ({
  filters,
  groupOptions,
  onGroupChange,
  onSearchChange,
  filterGroupLabel,
  filterAllLabel,
  searchPlaceholder,
}: DocumentFiltersProps) => {
  return (
    <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
      <FormControl fullWidth>
        <InputLabel>{filterGroupLabel}</InputLabel>
        <Select
          value={filters.selectedGroup}
          label={filterGroupLabel}
          onChange={(event) => onGroupChange(event.target.value)}
        >
          <MenuItem value="all">{filterAllLabel}</MenuItem>
          {groupOptions.map(option => (
            <MenuItem key={option.value} value={option.value}>
              {option.label}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      <TextField
        fullWidth
        value={filters.searchTerm}
        label={searchPlaceholder}
        onChange={(event) => onSearchChange(event.target.value)}
      />
    </Stack>
  );
};

export default DocumentFilters;
