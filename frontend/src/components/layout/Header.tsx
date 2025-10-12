import { useNavigate } from 'react-router-dom';
import { 
  AppBar, 
  Toolbar, 
  Typography, 
  IconButton, 
  Box,
  useTheme
} from '@mui/material';
import { Settings as SettingsIcon } from '@mui/icons-material';
import { useLanguage } from '../../contexts/LanguageContext';
import { ROUTES } from '../../constants';

const Header = () => {
  const navigate = useNavigate();
  const theme = useTheme();
  const { t } = useLanguage();

  const handleAccountSettings = () => {
    navigate(ROUTES.SETTINGS);
  };

  return (
    <AppBar 
      position="sticky" 
      elevation={0}
      sx={{
        backgroundColor: 'rgba(255, 255, 255, 0.95)',
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
            background: `linear-gradient(135deg, ${theme.palette.primary.main} 0%, ${theme.palette.secondary.main} 100%)`,
            backgroundClip: 'text',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            cursor: 'pointer',
          }}
        >
          MCD
        </Typography>
        <IconButton
          onClick={handleAccountSettings}
          title={t('header.settings')}
          sx={{
            backgroundColor: '#f8fafc',
            border: '1px solid #cbd5e1',
            color: '#475569',
            '&:hover': {
              backgroundColor: '#e2e8f0',
              borderColor: '#94a3b8',
              color: '#334155',
              transform: 'translateY(-2px)',
              boxShadow: '0 2px 6px rgba(0, 0, 0, 0.15)',
            },
            transition: 'all 0.3s ease',
          }}
        >
          <SettingsIcon />
        </IconButton>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
