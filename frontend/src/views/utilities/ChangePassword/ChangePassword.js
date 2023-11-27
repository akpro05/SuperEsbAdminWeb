import React, { useState, useEffect } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {Visibility,VisibilityOff} from '@mui/icons-material';
import { Dialog, DialogTitle, DialogContent, DialogActions  } from '@mui/material';
import { Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton,InputAdornment } from '@mui/material';

// Import the daterangepicker CSS
import axios from 'axios';
import { Formik, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';

// project imports
import MainCard from '../../../ui-component/cards/MainCard';
import Breadcrumbs from '@mui/material/Breadcrumbs';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';
import HomeIcon from '@mui/icons-material/Home';
import { CREATE_USER_URL, SEARCH_USER_URL } from '../../../UrlPath.js';

// styles
const IFrameWrapper = styled('iframe')(({ theme }) => ({
  height: 'calc(100vh - 210px)',
  border: '1px solid',
  borderColor: theme.palette.primary.light
}));

// =============================||Breadcrumb ||============================= //
function Breadcrumb({ label, onClick }) {
  return (
    <IconButton onClick={onClick}>
      {label}
    </IconButton>
  );
}

const TablerIcons = () => {
  const navigate = useNavigate();
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
  const [showOldPassword, setShowOldPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const { t, i18n } = useTranslation();

  const searchUserBreadcrumbs = [
    <Breadcrumbs aria-label="breadcrumb">
      <Link underline="hover" color="inherit" href="/Dashboard">
        <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
        {t('Home')}
      </Link>
      <Link underline="hover" color="inherit" href="#">
      {t('System User Management')}
      </Link>
      <Typography color="text.primary">{t('changePassword')} </Typography>
    </Breadcrumbs>
  ];

  // --------For status dropdown------
  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

 const [formData, setFormData] = useState({
     		  fullName: '',
              email: '',
              mobile: '',
              address: '',
              input_status: '',
  });

  const handleSubmit = async (values, { resetForm }) => {
    try {
      const response = await axios.post('/ChangePassword', values);
      if (response.data.success) {
        openPopupWithMessage('Change password created successfully');
        resetForm({});
      } else {
        setPopupMessage(response.data.message);
        setOpenPopup(true); // Open the popup dialog
        resetForm({});

      }
    } catch (error) {
      // Handle network or other errors
      console.error('Error creating user:', error);
      openPopupWithMessage('An error occurred');
      resetForm({});
    }
  };

// --------For Validation Message------//
 const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{6,10}$/;

  const validationSchema = Yup.object().shape({
  oldpassword: Yup.string()
    .required(t('oldPasswordRequired'))
    .min(6, t('minimumCharacters')),
  newpassword: Yup.string()
    .required(t('newPasswordRequired1'))
    .matches(passwordRegex, t('newPasswordRequired2'))
    .min(6, t('minimumCharacters'))
    .test('passwords-not-same', t('oldAndNewPasswordsSame'), function (value) {
      const oldPassword = this.parent.oldpassword;
      return oldPassword !== value;
    }),
  confirmpassword: Yup.string()
    .required(t('confirmPasswordRequired1'))
    .matches(passwordRegex, t('confirmPasswordRequired2'))
    .oneOf([Yup.ref('newpassword'), null], t('passwordsMustMatch'))
    .min(6, t('minimumCharacters')),
});



//-----------For Popup Messages--------//

const openPopupWithMessage = (message) => {
  setPopupMessage(message);
  setOpenPopup(true);
  };
  return (
    <>
      <MainCard title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('changePassword')}</span>} breadcrumbs={searchUserBreadcrumbs}>
        <Card sx={{ overflow: 'hidden' }}>
          <Formik
            initialValues={{
              oldpassword: '',
              newpassword: '',
              confirmpassword: '',
            }}
            validationSchema={validationSchema}
            onSubmit={handleSubmit}
          >
            {({ errors, touched, handleSubmit }) => (
              <form noValidate onSubmit={handleSubmit}>
                <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Grid item xs={6} style={{marginLeft:'70px'}}>
                    <Field
                      name="oldpassword"
                      as={TextField}
                      label={t('oldPasswordLabel')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                      type={showOldPassword ? 'text' : 'password'} // Toggle password visibility
					      InputProps={{
					        endAdornment: (
					          <InputAdornment position="end">
					            <IconButton
					              onClick={() => setShowOldPassword(!showOldPassword)} // Toggle show/hide password
					              edge="end"
					            >
					              {showOldPassword ? <Visibility /> : <VisibilityOff />}
					            </IconButton>
					          </InputAdornment>
					        ),
					      }}
                    />
                    {touched.oldpassword && errors.oldpassword && (
                      <FormHelperText error>{errors.oldpassword}</FormHelperText>
                    )}
                  </Grid>
                
                <Grid item xs={6} style={{marginLeft:'70px'}}>
                    <Field
                      name="newpassword"
                      as={TextField}
                      label={t('newPasswordLabel')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                    type={showNewPassword ? 'text' : 'password'} // Toggle password visibility
					      InputProps={{
					        endAdornment: (
					          <InputAdornment position="end">
					            <IconButton
					              onClick={() => setShowNewPassword(!showNewPassword)} // Toggle show/hide password
					              edge="end"
					            >
					              {showNewPassword ? <Visibility /> : <VisibilityOff />}
					            </IconButton>
					          </InputAdornment>
					        ),
					      }}
                    />
                    {touched.newpassword && errors.newpassword && (
                      <FormHelperText error>{errors.newpassword}</FormHelperText>
                    )}
                  </Grid>
                
                <Grid item xs={6} style={{marginLeft:'70px'}}>
                    <Field
                      name="confirmpassword"
                      as={TextField}
                      label={t('confirmPasswordLabel')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                    type={showConfirmPassword ? 'text' : 'password'} // Toggle password visibility
					      InputProps={{
					        endAdornment: (
					          <InputAdornment position="end">
					            <IconButton
					              onClick={() => setShowConfirmPassword(!showConfirmPassword)} // Toggle show/hide password
					              edge="end"
					            >
					              {showConfirmPassword ? <Visibility /> : <VisibilityOff />}
					            </IconButton>
					          </InputAdornment>
					        ),
					      }}
                    />
                    {touched.confirmpassword && errors.confirmpassword && (
                      <FormHelperText error>{errors.confirmpassword}</FormHelperText>
                    )}
                  </Grid>
                </Grid>
                <Grid item xs={6}>
				    <Box sx={{ p: 2 }}>
				      <Typography variant="body2" style={{fontSize:'16px',fontWeight:'bold'}}>
				        {t('passwordConditions')}
				        <br />
				        <br />
				        - {t('minimumCharacters')}
				        <br />
				        - {t('upperCases')}
				        <br />
				        - {t('specialCharacter')}
				        <br />
				        - {t('passwordsMustMatch')}
				        <br />
				        - {t('oldAndNewPasswordsSame')}
				      </Typography>
				    </Box>
				  </Grid>
                </Grid>
                  
                  
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
                  <Button variant="contained" color="primary" style={{ backgroundColor: '#1478c8' }} type="submit">
                    {t('submit')}
                  </Button>
                  <Box sx={{ mx: 2 }} />
                  <Button variant="contained" color="secondary" style={{ backgroundColor: 'rgb(13 32 97)' }} onClick={() => navigate('/Dashboard')}>
                    {t('back')}
                  </Button>
                </Box>
              </form>
            )}
          </Formik>
        </Card>
      </MainCard>
    <Dialog open={openPopup} onClose={() => setOpenPopup(false)}>
        <DialogTitle>{t('superEsbAdminAlert')}</DialogTitle>
        <DialogContent>
          {popupMessage}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenPopup(false)} color="primary">
            close
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default TablerIcons;
