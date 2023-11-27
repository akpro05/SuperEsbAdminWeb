import React from 'react';
import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogTitle, DialogContent, DialogActions  } from '@mui/material';
import axios from 'axios';
import { FORGOT_URL } from '../../../../UrlPath.js';

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
  TextField,
  Typography,
  useMediaQuery,
  Stack,
  Select,
  MenuItem
} from '@mui/material';

// third party
import * as Yup from 'yup';
import { Formik } from 'formik';

// project imports
import AuthWrapper1 from '../AuthWrapper1';
import AuthCardWrapper from '../AuthCardWrapper';
import Logo from '../../../../ui-component/Logo.js';
import AuthForgotPassword from '../auth-forms/AuthForgotPassword';
import AuthFooter from '../../../../ui-component/cards/AuthFooter';
import useScriptRef from '../../../../hooks/useScriptRef';
import Google from '../../../../assets/images/icons/social-google.svg';
import AnimateButton from '../../../../ui-component/extended/AnimateButton';
import { strengthColor, strengthIndicator } from '../../../../utils/password-strength';

// assets
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import '../../../../assets/scss/style.css'
import backgroundImage from '../../../../assets/images/esb2.jpg'

const containerStyles = {
  minHeight: '100vh',
  background: `url(${backgroundImage})`,
  backgroundBlendMode: 'overlay, overlay, normal', // You can also split this into an array if needed
};

// ===============================|| AUTH3 - REGISTER ||=============================== //

const Register = ({ ...others }) => {
  const theme = useTheme();
  const scriptedRef = useScriptRef();
  const matchDownSM = useMediaQuery(theme.breakpoints.down('md'));
  const customization = useSelector((state) => state.customization);
  const [showPassword, setShowPassword] = useState(false);
  const [checked, setChecked] = useState(true);
  const [email, setEmail] = useState('');
  const navigate = useNavigate();
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

const [strength, setStrength] = useState(0);
  const [level, setLevel] = useState();

  const googleHandler = async () => {
    console.error('Register');
  };

  const handleClickShowPassword = () => {
    setShowPassword(!showPassword);
  };

  const handleMouseDownPassword = (event) => {
    event.preventDefault();
  };

  const changePassword = (value) => {
    const temp = strengthIndicator(value);
    setStrength(temp);
    setLevel(strengthColor(temp));
  };

  useEffect(() => {
    changePassword('123456');
  }, []);

	const handleForgotPassword = async (values) => {
    try {
      const response = await axios.post(FORGOT_URL, {
        email: values.email,
        language: currentLanguage,
      });

      if (response.data.success) {
    // Password reset successful
    let successMessage = '';
    if (currentLanguage === 'english') {
        successMessage = 'Password reset successfully!';
    } else if (currentLanguage === 'french') {
        successMessage = 'Réinitialisation du mot de passe réussie !'; // French success message
    }
    setPopupMessage(successMessage);
    setOpenPopup(true);
    setEmail('');
} else {
    // Failed to reset password
    let failureMessage = '';
    if (currentLanguage === 'english') {
        failureMessage = 'Failed to reset password. Please try again.';
    } else if (currentLanguage === 'french') {
        failureMessage = 'Échec de la réinitialisation du mot de passe. Veuillez réessayer.'; // French failure message
    }
    setPopupMessage(failureMessage);
    setOpenPopup(true);
}
    } catch (error) {
    	let failureMessage = '';
    if (currentLanguage === 'english') {
        failureMessage = 'An error occurred. Please try again later.';
    } else if (currentLanguage === 'french') {
        failureMessage = "Une erreur s'est produite. Veuillez réessayer plus tard."; // French failure message
    }
      setPopupMessage(failureMessage);
      setOpenPopup(true);
    }
  };
const handleClosePopup = () => {
    setOpenPopup(false);
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
                            {t('forgotPassword1')} 
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
      </Grid>

      <Formik
        initialValues={{
          email: '',
          submit: null
        }}
        validationSchema={Yup.object().shape({
          email: Yup.string().email('Must be a valid email').max(255).required(t('Emailrequired')),
        })}
        onSubmit={handleForgotPassword}
      >
        {({ errors, handleBlur, handleChange, handleSubmit, isSubmitting, touched, values }) => (
          <form noValidate onSubmit={handleSubmit} {...others}>
            <FormControl fullWidth error={Boolean(touched.email && errors.email)} sx={{ ...theme.typography.customInput }}>
              <InputLabel htmlFor="outlined-adornment-email-register">{t('yourName')}</InputLabel>
              <OutlinedInput
                id="outlined-adornment-email-register"
                type="email"
                value={values.email}
                name="email"
                onBlur={handleBlur}
                onChange={handleChange}
                inputProps={{}}
              />
              {touched.email && errors.email && (
                <FormHelperText error id="standard-weight-helper-text--register">
                  {errors.email}
                </FormHelperText>
              )}
            </FormControl>
            <Box sx={{ mt: 2 }}>
              <AnimateButton>
                <Button disableElevation disabled={isSubmitting} style={{ backgroundColor: 'rgb(13 32 97)' }} fullWidth size="large" type="submit" variant="contained" color="secondary">
                 {t('submit')}
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
                  <Grid item xs={12}>
                    <Grid item container direction="column" alignItems="center" xs={12}>
                      <Typography component={Link} to="/Login" variant="subtitle1" sx={{ textDecoration: 'none' }}>
                        {t('Already')}
                      </Typography>
                    </Grid>
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
        <DialogTitle>Message</DialogTitle>
        <DialogContent>
          {popupMessage}
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

export default Register;