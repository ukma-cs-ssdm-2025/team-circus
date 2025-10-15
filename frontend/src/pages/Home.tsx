import { Box, Stack } from '@mui/material';
import { Sidebar } from '../components/layout';
import { ActionButton } from '../components/forms';
import { useLanguage } from '../contexts/LanguageContext';
import type { BaseComponentProps } from '../types';

interface HomeProps extends BaseComponentProps {}

const Home = ({ className = '' }: HomeProps) => {
  const { t } = useLanguage();

  return (
    <Box 
      className={className}
      sx={{
        display: 'flex',
        gap: 2,
        minHeight: '100vh',
        px: 2,
        py: 3,
      }}
    >
      <Box sx={{ display: { xs: 'none', md: 'block' } }}>
        <Sidebar />
      </Box>

      <Box
        sx={{
          flex: 1,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Stack
          direction={{ xs: 'column', sm: 'row' }}
          spacing={3}
          justifyContent="center"
          alignItems="center"
        >
          <ActionButton
            onClick={() => alert(t('home.createDocumentAlert'))}
          >
            {t('home.createDocument')}
          </ActionButton>
          <ActionButton
            onClick={() => alert(t('home.createGroupAlert'))}
          >
            {t('home.createGroup')}
          </ActionButton>
        </Stack>
      </Box>
    </Box>
  );
};

export default Home;
