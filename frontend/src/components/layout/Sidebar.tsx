import {
  Box,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Typography,
  Divider,
  useTheme,
} from '@mui/material';
import {
  Group as GroupIcon,
  Home as HomeIcon,
  Description as DescriptionIcon,
  Settings as SettingsIcon,
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useLanguage } from '../../contexts/LanguageContext';
import type { BaseComponentProps } from '../../types';
import { ROUTES } from '../../constants';

type SidebarProps = BaseComponentProps;

const Sidebar = ({ className = '' }: SidebarProps) => {
  const theme = useTheme();
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useLanguage();

  const isActive = (path: string) => location.pathname === path;

  return (
    <Box
      className={className}
      component="nav"
      sx={{
        width: 280,
        flexShrink: 0,
        p: 2,
        position: { xs: 'static', md: 'sticky' },
        top: { md: 88 },
        alignSelf: 'flex-start',
        backgroundColor: theme.palette.mode === 'light'
          ? 'rgba(255, 255, 255, 0.75)'
          : 'rgba(255, 255, 255, 0.08)',
        borderRight: `1px solid ${theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.2)' : theme.palette.divider}`,
        borderRadius: 2,
        backdropFilter: 'blur(8px)',
      }}
    >
      <Typography variant="h6" sx={{ px: 1, pb: 1, fontWeight: 700 }}>
        {t('sidebar.navigation')}
      </Typography>
      <List>
        <ListItemButton selected={isActive(ROUTES.HOME)} onClick={() => navigate(ROUTES.HOME)}>
          <ListItemIcon><HomeIcon /></ListItemIcon>
          <ListItemText primary={t('sidebar.home')} />
        </ListItemButton>
        <ListItemButton selected={isActive(ROUTES.DOCUMENTS)} onClick={() => navigate(ROUTES.DOCUMENTS)}>
          <ListItemIcon><DescriptionIcon /></ListItemIcon>
          <ListItemText primary={t('sidebar.documents')} />
        </ListItemButton>
        <ListItemButton selected={isActive(ROUTES.SETTINGS)} onClick={() => navigate(ROUTES.SETTINGS)}>
          <ListItemIcon><SettingsIcon /></ListItemIcon>
          <ListItemText primary={t('sidebar.settings')} />
        </ListItemButton>
      </List>

      <Divider sx={{ my: 1 }} />

      <Typography variant="subtitle2" sx={{ px: 1, py: 1, fontWeight: 600, opacity: 0.8 }}>
        {t('sidebar.groups')}
      </Typography>
      <List>
        <ListItemButton onClick={() => navigate(ROUTES.GROUPS)}>
          <ListItemIcon><GroupIcon /></ListItemIcon>
          <ListItemText primary={t('sidebar.viewGroups')} />
        </ListItemButton>
        <ListItemButton onClick={() => alert(t('home.createGroupAlert'))}>
          <ListItemIcon><GroupIcon /></ListItemIcon>
          <ListItemText primary={t('sidebar.createGroup')} />
        </ListItemButton>
      </List>
    </Box>
  );
};

export default Sidebar;
