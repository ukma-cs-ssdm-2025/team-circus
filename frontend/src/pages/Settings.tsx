import { useState } from 'react';
import {
  Box,
  Container,
  Typography,
  Card,
  CardContent,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormControlLabel,
  Switch,
  Button,
  Stack,
  useTheme as useMuiTheme
} from '@mui/material';
import { useLanguage } from '../contexts/LanguageContext';
import { useTheme } from '../contexts/ThemeContext';
import type { BaseComponentProps } from '../types';

interface SettingsProps extends BaseComponentProps {}

const Settings = ({ className = '' }: SettingsProps) => {
  const muiTheme = useMuiTheme();
  const { theme, setTheme } = useTheme();
  const { language, setLanguage, t } = useLanguage();
  const [settings, setSettings] = useState({
    notifications: true,
    autoSave: true,
  });

  const handleSettingChange = (key: string, value: any) => {
    if (key === 'language') {
      setLanguage(value);
    } else if (key === 'theme') {
      setTheme(value);
    } else {
      setSettings(prev => ({
        ...prev,
        [key]: value
      }));
    }
  };

  return (
    <Box className={className}>
      <Container maxWidth="md" sx={{ py: 8 }}>
        <Card
          sx={{
            background: muiTheme.palette.mode === 'light' 
              ? 'rgba(255, 255, 255, 0.8)' 
              : 'rgba(30, 30, 30, 0.8)',
            backdropFilter: 'blur(10px)',
            borderRadius: 4,
            boxShadow: '0 10px 30px rgba(0, 0, 0, 0.1)',
          }}
        >
          <CardContent sx={{ p: 6 }}>
                   <Typography
                     variant="h3"
                     sx={{
                       fontWeight: 700,
                       color: muiTheme.palette.primary.main,
                       mb: 1,
                     }}
                   >
                     {t('settings.title')}
                   </Typography>
                   
                   <Typography variant="h6" color="text.secondary" sx={{ mb: 4 }}>
                     {t('settings.subtitle')}
                   </Typography>

            <Stack spacing={4}>
              <Card variant="outlined" sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom>
                  {t('settings.general')}
                </Typography>
                
                <Stack spacing={3}>
                         <FormControl fullWidth>
                           <InputLabel>{t('settings.theme')}</InputLabel>
                           <Select
                             value={theme}
                             label={t('settings.theme')}
                             onChange={(e) => handleSettingChange('theme', e.target.value)}
                           >
                             <MenuItem value="light">{t('settings.theme.light')}</MenuItem>
                             <MenuItem value="dark">{t('settings.theme.dark')}</MenuItem>
                           </Select>
                         </FormControl>

                  <FormControl fullWidth>
                    <InputLabel>{t('settings.language')}</InputLabel>
                    <Select
                      value={language}
                      label={t('settings.language')}
                      onChange={(e) => handleSettingChange('language', e.target.value)}
                    >
                      <MenuItem value="uk">Українська</MenuItem>
                      <MenuItem value="en">English</MenuItem>
                    </Select>
                  </FormControl>
                </Stack>
              </Card>

              <Card variant="outlined" sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom>
                  {t('settings.notifications')}
                </Typography>
                
                <FormControlLabel
                  control={
                    <Switch
                      checked={settings.notifications}
                      onChange={(e) => handleSettingChange('notifications', e.target.checked)}
                    />
                  }
                  label={t('settings.notificationsLabel')}
                />
              </Card>

              <Card variant="outlined" sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom>
                  {t('settings.documents')}
                </Typography>
                
                <FormControlLabel
                  control={
                    <Switch
                      checked={settings.autoSave}
                      onChange={(e) => handleSettingChange('autoSave', e.target.checked)}
                    />
                  }
                  label={t('settings.autoSaveLabel')}
                />
              </Card>
            </Stack>

            <Stack
              direction={{ xs: 'column', sm: 'row' }}
              spacing={2}
              justifyContent="center"
              sx={{ mt: 4 }}
            >
              <Button
                variant="contained"
                size="large"
                sx={{ minWidth: 200 }}
              >
                {t('settings.save')}
              </Button>
              <Button
                variant="outlined"
                size="large"
                sx={{ minWidth: 200 }}
              >
                {t('settings.reset')}
              </Button>
            </Stack>
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
};

export default Settings;
