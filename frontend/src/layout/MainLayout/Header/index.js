import React, { useState ,useEffect} from 'react';
import PropTypes from 'prop-types';

// material-ui
import { useTheme } from '@mui/material/styles';
import { Avatar, Box, ButtonBase,IconButton } from '@mui/material';
import { useTranslation } from 'react-i18next';

// project imports
import LogoSection from '../LogoSection';
import SearchSection from './SearchSection';
import ProfileSection from './ProfileSection';
import NotificationSection from './NotificationSection';

// assets
import { IconMenu2 } from '@tabler/icons';
import enFlagImage from '../../../assets/images/en-flag.png';
import frFlagImage from '../../../assets/images/fn-flag.png';

// ==============================|| MAIN NAVBAR / HEADER ||============================== //

const Header = ({ handleLeftDrawerToggle }) => {
  const theme = useTheme();
  const [currentLanguage, setCurrentLanguage] = useState('english');
  const { i18n } = useTranslation();

  const toggleLanguage = () => {
    const newLanguage = currentLanguage === 'english' ? 'french' : 'english';
    setCurrentLanguage(newLanguage);
    i18n.changeLanguage(newLanguage);
  };

  // When the component mounts, check the language from the session storage
  useEffect(() => {
    const sessionLanguage = sessionStorage.getItem('language');
    if (sessionLanguage) {
      setCurrentLanguage(sessionLanguage);
      i18n.changeLanguage(sessionLanguage);
    }
  }, [i18n]);

  return (
    <>
      {/* logo & toggler button */}
      <Box
        sx={{
          width: 228,
          display: 'flex',
          [theme.breakpoints.down('md')]: {
            width: 'auto'
          }
        }}
      >
        <Box component="span" sx={{ display: { xs: 'none', md: 'block' }, flexGrow: 1 }}>
          <LogoSection />
        </Box>
        <ButtonBase sx={{ borderRadius: '12px', overflow: 'hidden' }}>
          <Avatar
            variant="rounded"
            sx={{
              ...theme.typography.commonAvatar,
              ...theme.typography.mediumAvatar,
              transition: 'all .2s ease-in-out',
              background: '#03266a4a',
              color: 'rgb(13 32 97)',
              '&:hover': {
                background: 'rgb(13 32 97)',
                color: 'white'
              }
            }}
            onClick={handleLeftDrawerToggle}
            color="inherit"
          >
            <IconMenu2 stroke={1.5} size="1.3rem" />
          </Avatar>
        </ButtonBase>
      </Box>

      {/* header search */}
     
      <Box sx={{ flexGrow: 1 }} />
      <Box sx={{ flexGrow: 1 }} />
		{currentLanguage !== sessionStorage.getItem('language') && (
        <IconButton
          size="large"
          aria-label={
            currentLanguage === 'english'
              ? 'Switch to French'
              : 'Switch to English'
          }
          onClick={toggleLanguage}
          color="inherit"
        >
          {currentLanguage === 'english' ? (
            <img src={enFlagImage} alt="English Flag" width={24} height={24} />
          ) : (
            <img src={frFlagImage} alt="French Flag" width={24} height={24} />
          )}
        </IconButton>
      )}
      {/* notification & profile */}
     {/* <NotificationSection />*/}
      <ProfileSection />
    </>
  );
};

Header.propTypes = {
  handleLeftDrawerToggle: PropTypes.func
};

export default Header;