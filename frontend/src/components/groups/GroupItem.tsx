import { Box, Divider, ListItem, ListItemButton, ListItemText, Typography } from '@mui/material';
import { formatDate } from '../../utils';
import type { GroupItem as GroupItemType } from '../../types/entities';

interface GroupItemProps {
  group: GroupItemType;
  isLast: boolean;
  onClick: (groupUUID: string) => void;
  createdAtLabel: string;
}

const GroupItem = ({ group, isLast, onClick, createdAtLabel }: GroupItemProps) => {
  return (
    <Box>
      <ListItem disablePadding alignItems="flex-start">
        <ListItemButton
          sx={{ borderRadius: 2 }}
          onClick={() => onClick(group.uuid)}
        >
          <ListItemText
            primary={
              <Typography variant="h6" sx={{ fontWeight: 600 }}>
                {group.name}
              </Typography>
            }
            secondary={
              <Typography variant="body2" color="text.secondary">
                {`${createdAtLabel}: ${formatDate(group.created_at)}`}
              </Typography>
            }
          />
        </ListItemButton>
      </ListItem>
      {!isLast && <Divider sx={{ my: 1 }} />}
    </Box>
  );
};

export default GroupItem;
