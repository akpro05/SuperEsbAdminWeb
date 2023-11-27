import React from 'react';
import { useState,useEffect } from 'react';
import { useSelector } from 'react-redux';
import { Link } from 'react-router-dom';

// material-ui
import { useTheme } from '@mui/material/styles';
import {
  Box,
  Button,
  Checkbox,
  Divider,
  FormControl,
  FormControlLabel,
  FormHelperText,
  Grid,
  IconButton,
  InputAdornment,
  InputLabel,
  OutlinedInput,
  Stack,
  Typography,
  useMediaQuery,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Select,
  MenuItem
} from '@mui/material';
import { useTranslation } from 'react-i18next';

// third party
import * as Yup from 'yup';
import { Formik } from 'formik';

// project imports
import useScriptRef from '../../../../hooks/useScriptRef';
import AnimateButton from '../../../../ui-component/extended/AnimateButton';

// project imports
import AuthWrapper1 from '../AuthWrapper1';
import AuthCardWrapper from '../AuthCardWrapper';
import AuthLogin from '../auth-forms/AuthLogin';
import Logo from '../../../../ui-component/Logo';
import AuthFooter from '../../../../ui-component/cards/AuthFooter';

// assets
import axios from 'axios';
import '../../../../assets/scss/style.css'
import backgroundImage from '../../../../assets/images/esb2.jpg'
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import '../../../../assets/scss/style.css'
import Google from '../../../../assets/images/icons/social-google.svg';

// ================================|| AUTH3 - LOGIN ||=============================== //

const Login = ({ ...others }) => {
  const theme = useTheme();
  const scriptedRef = useScriptRef();
  const matchDownSM = useMediaQuery(theme.breakpoints.down('md'));
  const customization = useSelector((state) => state.customization);
  const [checked, setChecked] = useState(true);
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
  const [currentLanguage, setCurrentLanguage] = useState('english');
  const { t, i18n } = useTranslation();

const toggleLanguage = async () => {
  const newLanguage = currentLanguage === 'english' ? 'french' : 'english';
  setCurrentLanguage(newLanguage);
  i18n.changeLanguage(newLanguage);

  const requestData = { language: newLanguage };

  console.log('Request Data:', requestData); // Log the request data

  try {
    // Make an API request to update the language on the backend
    await axios.post('/Login', requestData);
    // Replace '/your-backend-endpoint-for-updating-language' with the actual endpoint on your backend
  } catch (error) {
    console.error('Error updating language:', error);
    // Handle error appropriately (e.g., show a notification)
  }
};

  const googleHandler = async () => {
    console.error('Login');
  };

  const [showPassword, setShowPassword] = useState(false);
  const handleClickShowPassword = () => {
    setShowPassword(!showPassword);
  };

  const handleMouseDownPassword = (event) => {
    event.preventDefault();
  };

  const handleClosePopup = () => {
    setOpenPopup(false);
  };

  const handleLogin = async (values) => {
  try {
    const response = await axios.post('/Login', {
      email: values.email,
      password: values.password,
      language: currentLanguage, // Include the language property
    });

    if (response.data.status) {
      // Successful login, redirect to the dashboard
      window.location.href = '/Dashboard';
      sessionStorage.setItem('menu', response.data.menu);
      sessionStorage.setItem('submenu', response.data.submenu);
      sessionStorage.setItem('language', response.data.language);
    } else {
      // Unsuccessful login, set the error message and open the popup
      setPopupMessage(response.data.message);
      setOpenPopup(true);
    }
  } catch (error) {
    console.error('Error during login:', error);
  }
};

const validationSchema = Yup.object().shape({
  email: Yup.string()
    .email(t('invalidEmail')) // Translate the error message
    .max(255, t('EmailmaxLength'))
    .required(t('Emailrequired')),
  password: Yup.string()
    .max(255, t('PasswordmaxLength'))
    .required(t('Passwordrequired')),
});



  const containerStyles = {
    minHeight: '100vh',
    background: `url(${backgroundImage})`,
    backgroundBlendMode: 'overlay, overlay, normal',
  };

  return (
    <AuthWrapper1>
      <Grid container direction="column" justifyContent="flex-end" sx={containerStyles}>
        <Select
          className="btn-block"
          value={currentLanguage}
          onChange={toggleLanguage}
        >
          <MenuItem value="english">English</MenuItem>
          <MenuItem value="french">French</MenuItem>
        </Select>
        <Grid item xs={12}>
          <Grid container justifyContent="center" alignItems="center" sx={{ minHeight: 'calc(100vh - 68px)' }}>
            <Grid item sx={{ m: { xs: 1, sm: 3 }, mb: 0 }}>
              <AuthCardWrapper>
                <Grid container spacing={2} alignItems="center" justifyContent="center">
                  <Grid item sx={{ mb: 3 }}>
                    <Link to="#">
                      <Logo />
                    </Link>
                  </Grid>
                  <Grid item xs={12}>
                    <Grid container direction={matchDownSM ? 'column-reverse' : 'row'} alignItems="center" justifyContent="center">
                      <Grid item>
                        <Stack alignItems="center" justifyContent="center" spacing={1}>
                          <Typography color={theme.palette.secondary.main} style={{ color: 'rgb(13 32 97)' }} gutterBottom variant={matchDownSM ? 'h3' : 'h2'}>
                            {t('welcome')}
                          </Typography>
                          <Typography color={theme.palette.secondary.main} style={{ color: 'rgb(13 32 97)' }} gutterBottom variant={matchDownSM ? 'h3' : 'h2'}>
                            {t('superESB')}
                          </Typography>
                          <Typography variant="caption" fontSize="16px" textAlign={matchDownSM ? 'center' : 'inherit'}>
                            {t('credentials')}
                          </Typography>
                        </Stack>
                      </Grid>
                    </Grid>
                  </Grid>
                  <Grid item xs={12}>
                    <Grid container direction="column" justifyContent="center" spacing={2}>
        
        <Grid item xs={12}>
          <Box
            sx={{
              alignItems: 'center',
              display: 'flex'
            }}
          >
            <Divider sx={{ flexGrow: 1 }} orientation="horizontal" />

          

            <Divider sx={{ flexGrow: 1 }} orientation="horizontal" />
          </Box>
        </Grid>
        <Grid item xs={12} container alignItems="center" justifyContent="center">
          <Box sx={{ mb: 2 }}>
            <Typography variant="subtitle1">{t('signInWith')}</Typography>
          </Box>
        </Grid>
      </Grid>
    
    <Formik
  initialValues={{
    email: '',
    password: '',
  }}
  validationSchema={validationSchema}
  onSubmit={handleLogin} 
>
        {({ errors, handleBlur, handleChange, handleSubmit, isSubmitting, touched, values }) => (
          <form noValidate onSubmit={handleSubmit} {...others}>
            <FormControl fullWidth error={Boolean(touched.email && errors.email)} sx={{ ...theme.typography.customInput }}>
              <InputLabel htmlFor="outlined-adornment-email-login">{t('yourName')}</InputLabel>
              <OutlinedInput
                id="outlined-adornment-email-login"
                type="email"
                value={values.email}
                name="email"
                onBlur={handleBlur}
                onChange={handleChange}
                label="Email Address / Username"
                inputProps={{}}
              />
              {touched.email && errors.email && (
                <FormHelperText error id="standard-weight-helper-text-email-login">
                  {errors.email}
                </FormHelperText>
              )}
            </FormControl>

            <FormControl fullWidth error={Boolean(touched.password && errors.password)} sx={{ ...theme.typography.customInput }}>
              <InputLabel htmlFor="outlined-adornment-password-login">{t('password')}</InputLabel>
              <OutlinedInput
                id="outlined-adornment-password-login"
                type={showPassword ? 'text' : 'password'}
                value={values.password}
                name="password"
                onBlur={handleBlur}
                onChange={handleChange}
                endAdornment={
                  <InputAdornment position="end">
                    <IconButton
                      aria-label="toggle password visibility"
                      onClick={handleClickShowPassword}
                      onMouseDown={handleMouseDownPassword}
                      edge="end"
                      size="large"
                    >
                      {showPassword ? <Visibility /> : <VisibilityOff />}
                    </IconButton>
                  </InputAdornment>
                }
                label="Password"
                inputProps={{}}
              />
              {touched.password && errors.password && (
                <FormHelperText error id="standard-weight-helper-text-password-login">
                  {errors.password}
                </FormHelperText>
              )}
            </FormControl>
            <Stack direction="row" alignItems="center" justifyContent="space-between" spacing={1}>
            
              <Typography variant="subtitle1" color="secondary" sx={{ textDecoration: 'none', cursor: 'pointer',color: 'rgb(13 32 97)' }} component={Link} to="/ForgotPassword">
                {t('forgotPassword')}
              </Typography>
            </Stack>
            {errors.submit && (
              <Box sx={{ mt: 3 }}>
                <FormHelperText error>{errors.submit}</FormHelperText>
              </Box>
            )}

            <Box sx={{ mt: 2 }}>
              <AnimateButton>
                <Button disableElevation disabled={isSubmitting} style={{ backgroundColor: 'rgb(13 32 97)' }} fullWidth size="large" type="submit" variant="contained" color="secondary">
                  {t('signIn')}
                </Button>
              </AnimateButton>
            </Box>
          </form>
        )}
      </Formik>
                  </Grid>
                  <Grid item xs={12}>
                    <Divider />
                  </Grid>
                </Grid>
              </AuthCardWrapper>
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12} sx={{ m: 3, mt: 1 }}>
          <AuthFooter />
        </Grid>
      </Grid>
    
    <Dialog open={openPopup} onClose={handleClosePopup}>
        <DialogTitle>Error</DialogTitle>
        <DialogContent>
          <div className="error-message">{popupMessage}</div>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClosePopup} color="primary">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </AuthWrapper1>
  );
};

export default Login;