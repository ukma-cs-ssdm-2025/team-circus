import { useNavigate } from 'react-router-dom';
import { 
  AppBar, 
  Toolbar, 
  Typography, 
  IconButton, 
  Box,
  useTheme
} from '@mui/material';
import { Settings as SettingsIcon, LightMode as LightModeIcon, DarkMode as DarkModeIcon } from '@mui/icons-material';
import { useLanguage } from '../../contexts/LanguageContext';
import { useTheme as useAppTheme } from '../../contexts/ThemeContext';
import { ROUTES } from '../../constants';

const Header = () => {
  const navigate = useNavigate();
  const theme = useTheme();
  const { t } = useLanguage();
  const { theme: appTheme, toggleTheme } = useAppTheme();

  const handleAccountSettings = () => {
    navigate(ROUTES.SETTINGS);
  };

  return (
    <AppBar 
      position="sticky" 
      elevation={0}
      sx={{
        backgroundColor: theme.palette.mode === 'light' 
          ? 'rgba(255, 255, 255, 0.95)' 
          : 'rgba(30, 30, 30, 0.95)',
        backdropFilter: 'blur(10px)',
        boxShadow: '0 2px 20px rgba(0, 0, 0, 0.1)',
      }}
    >
      <Toolbar>
        <Typography
          variant="h5"
          component="a"
          href={ROUTES.HOME}
          sx={{
            flexGrow: 1,
            fontWeight: 700,
            textDecoration: 'none',
            color: theme.palette.primary.main,
            cursor: 'pointer',
          }}
        >
          MCD
        </Typography>
        
        <Box sx={{ display: 'flex', gap: 1 }}>
          <IconButton
            onClick={toggleTheme}
            title={t('header.toggleTheme')}
            sx={{
              backgroundColor: theme.palette.mode === 'light' ? '#f8fafc' : '#374151',
              border: `1px solid ${theme.palette.mode === 'light' ? '#cbd5e1' : '#4b5563'}`,
              color: theme.palette.mode === 'light' ? '#475569' : '#d1d5db',
              '&:hover': {
                backgroundColor: theme.palette.mode === 'light' ? '#e2e8f0' : '#4b5563',
                borderColor: theme.palette.mode === 'light' ? '#94a3b8' : '#6b7280',
                color: theme.palette.mode === 'light' ? '#334155' : '#f3f4f6',
                transform: 'translateY(-2px)',
                boxShadow: '0 2px 6px rgba(0, 0, 0, 0.15)',
              },
              transition: 'all 0.3s ease',
            }}
          >
            {appTheme === 'light' ? <DarkModeIcon /> : <LightModeIcon />}
          </IconButton>
          
          <IconButton
            onClick={handleAccountSettings}
            title={t('header.settings')}
            sx={{
              backgroundColor: theme.palette.mode === 'light' ? '#f8fafc' : '#374151',
              border: `1px solid ${theme.palette.mode === 'light' ? '#cbd5e1' : '#4b5563'}`,
              color: theme.palette.mode === 'light' ? '#475569' : '#d1d5db',
              '&:hover': {
                backgroundColor: theme.palette.mode === 'light' ? '#e2e8f0' : '#4b5563',
                borderColor: theme.palette.mode === 'light' ? '#94a3b8' : '#6b7280',
                color: theme.palette.mode === 'light' ? '#334155' : '#f3f4f6',
                transform: 'translateY(-2px)',
                boxShadow: '0 2px 6px rgba(0, 0, 0, 0.15)',
              },
              transition: 'all 0.3s ease',
            }}
          >
            <SettingsIcon />
          </IconButton>
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
