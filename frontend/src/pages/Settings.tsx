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
  useTheme
} from '@mui/material';
import type { BaseComponentProps } from '../types';

interface SettingsProps extends BaseComponentProps {}

const Settings = ({ className = '' }: SettingsProps) => {
  const theme = useTheme();
  const [settings, setSettings] = useState({
    theme: 'light',
    language: 'uk',
    notifications: true,
    autoSave: true,
  });

  const handleSettingChange = (key: string, value: any) => {
    setSettings(prev => ({
      ...prev,
      [key]: value
    }));
  };

  return (
    <Box className={className}>
      <Container maxWidth="md" sx={{ py: 8 }}>
        <Card
          sx={{
            background: 'rgba(255, 255, 255, 0.8)',
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
                background: `linear-gradient(135deg, ${theme.palette.primary.main} 0%, ${theme.palette.secondary.main} 100%)`,
                backgroundClip: 'text',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                mb: 1,
              }}
            >
              Налаштування акаунту
            </Typography>
            
            <Typography variant="h6" color="text.secondary" sx={{ mb: 4 }}>
              Керуйте своїми налаштуваннями та преференціями
            </Typography>

            <Stack spacing={4}>
              <Card variant="outlined" sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom>
                  Загальні налаштування
                </Typography>
                
                <Stack spacing={3}>
                  <FormControl fullWidth>
                    <InputLabel>Тема</InputLabel>
                    <Select
                      value={settings.theme}
                      label="Тема"
                      onChange={(e) => handleSettingChange('theme', e.target.value)}
                    >
                      <MenuItem value="light">Світла</MenuItem>
                      <MenuItem value="dark">Темна</MenuItem>
                      <MenuItem value="auto">Автоматично</MenuItem>
                    </Select>
                  </FormControl>

                  <FormControl fullWidth>
                    <InputLabel>Мова</InputLabel>
                    <Select
                      value={settings.language}
                      label="Мова"
                      onChange={(e) => handleSettingChange('language', e.target.value)}
                    >
                      <MenuItem value="uk">Українська</MenuItem>
                      <MenuItem value="en">English</MenuItem>
                      <MenuItem value="ru">Русский</MenuItem>
                    </Select>
                  </FormControl>
                </Stack>
              </Card>

              <Card variant="outlined" sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom>
                  Сповіщення
                </Typography>
                
                <FormControlLabel
                  control={
                    <Switch
                      checked={settings.notifications}
                      onChange={(e) => handleSettingChange('notifications', e.target.checked)}
                    />
                  }
                  label="Отримувати сповіщення"
                />
              </Card>

              <Card variant="outlined" sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom>
                  Робота з документами
                </Typography>
                
                <FormControlLabel
                  control={
                    <Switch
                      checked={settings.autoSave}
                      onChange={(e) => handleSettingChange('autoSave', e.target.checked)}
                    />
                  }
                  label="Автоматичне збереження"
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
                Зберегти зміни
              </Button>
              <Button
                variant="outlined"
                size="large"
                sx={{ minWidth: 200 }}
              >
                Скинути до стандартних
              </Button>
            </Stack>
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
};

export default Settings;
