import { List, Stack, Typography } from '@mui/material';
import GroupItem from './GroupItem';
import type { GroupItem as GroupItemType } from '../../types/entities';

interface GroupsListProps {
  groups: GroupItemType[];
  onGroupClick: (groupUUID: string) => void;
  totalLabel: string;
  createdAtLabel: string;
}

const GroupsList = ({ groups, onGroupClick, totalLabel, createdAtLabel }: GroupsListProps) => {
  return (
    <Stack spacing={2}>
      <Typography variant="body2" color="text.secondary">
        {`${totalLabel} ${groups.length}`}
      </Typography>

      <List disablePadding>
        {groups.map((group, index) => (
          <GroupItem
            key={group.uuid}
            group={group}
            isLast={index === groups.length - 1}
            onClick={onGroupClick}
            createdAtLabel={createdAtLabel}
          />
        ))}
      </List>
    </Stack>
  );
};

export default GroupsList;
