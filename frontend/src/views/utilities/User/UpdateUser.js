import React, { useState, useEffect } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import { Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton } from '@mui/material';

// Import the daterangepicker CSS
import axios from 'axios';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';

// project imports
import MainCard from '../../../ui-component/cards/MainCard';
import HomeIcon from '@mui/icons-material/Home';
import { CREATE_USER_URL, SEARCH_USER_URL } from '../../../UrlPath.js';
import Breadcrumbs from '@mui/material/Breadcrumbs';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';

// styles
const IFrameWrapper = styled('iframe')(({ theme }) => ({
  height: 'calc(100vh - 210px)',
  border: '1px solid',
  borderColor: theme.palette.primary.light
}));

// =============================|| Breadcrumb ||============================= //
function Breadcrumb({ label, onClick }) {
  return (
    <IconButton onClick={onClick}>
      {label}
    </IconButton>
  );
}

const TablerIcons = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [userData, setUserData] = useState(null);
  const [currentLanguage, setCurrentLanguage] = useState('en');
  const { t, i18n } = useTranslation();
  const [data, setData] = useState(null);

  const searchUserBreadcrumbs = [
    <Breadcrumbs aria-label="breadcrumb">
      <Link underline="hover" color="inherit" href="/Dashboard">
        <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
        {t('Home')}
      </Link>
      <Link underline="hover" color="inherit" href="#">
      {t('System User Management')}
      </Link>
      <Typography color="text.primary">{t('updateUser')}</Typography>
    </Breadcrumbs>
  ];
  

  // --------For status dropdown------
  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

  const hardcodedLanguages = [
    { id: 'english', name: t('english'), translation: t('english') },
    { id: 'french', name: t('french'), translation: t('french') },
];


  const initialValues = {
    fullName: '',
    email: '',
    mobile: '',
    address: '',
    input_status: '',
    input_role: '',
  };

  const validationSchema = Yup.object().shape({
  fullName: Yup.string()
    .matches(/^[A-Za-z\s]+$/, t('validationFullName1'))
    .max(30, t('validationFullName2'))
    .required(t('validationFullName3')),

  email: Yup.string()
    .max(30, t('validationEmail1'))
    .email(t('validationEmail2'))
    .required(t('validationEmail3')),
  mobile: Yup.string()
    .matches(/^\d{10}$/, t('validationMobile1'))
    .required(t('validationMobile2')),
  address: Yup.string()
  .max(50, t('validationAddress1'))
  .required(t('validationAddress2')),
  input_role: Yup.string().required(t('validationRolename')),
  input_status: Yup.string().required(t('validationStatus')),
  input_language: Yup.string().required(t('validationLanguage')),
});

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await axios.get(`/SysUser/UpdatesysUser/${id}`);
        setUserData(response.data);
        setData(response.data);
        setLoading(false);
        console.log(response.data);
      } catch (error) {
        setError('Failed to fetch data: ' + error.message);
        setLoading(false);
      }
    };

    fetchUserData();
  }, [id]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  const handleSubmit = async (values, { setValues }) => {

    try {
      const response = await axios.post(`/SysUser/UpdatesysUser/${id}`, values);

      if (response.data.message) {
        setPopupMessage(response.data.message);
        setOpenPopup(true); // Open the popup dialog

        // setValues({
        //   fullName: '',
        //   email: '',
        //   mobile: '',
        //   address: '',
        //   input_status: '',
        // }); // Reset the form fields to empty values
      } else {
        openPopupWithMessage('An error occurred');
        setOpenPopup(true); // Open the popup dialog
      }
    } catch (error) {
      // Handle network or other errors
      console.error('Error updating user:', error);
      openPopupWithMessage('An error occurred');

      // setValues({
      //   fullName: '',
      //   email: '',
      //   mobile: '',
      //   address: '',
      //   input_status: '',
      // }); // Reset the form fields to empty values
    }
  };




  // -----------For Popup Messages--------//

  const openPopupWithMessage = (message) => {
    setPopupMessage(message);
    setOpenPopup(true);
  };

  const handleDialogClose = () => {
    setOpenPopup(false);
    navigate('/SysUser/SearchsysUser'); // Redirect to /SysUser/SearchsysUser
  };


  return (
    <>
      <MainCard  title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('updateUser')}</span>} breadcrumbs={searchUserBreadcrumbs}>
        <Card sx={{ overflow: 'hidden' }}>
          <Formik
            initialValues={{
              fullName: userData?.FullName || '', // Populate the value from userData
              email: userData?.Email || '', // Populate the value from userData
              mobile: userData?.Mobile || '', // Populate the value from userData
              address: userData?.Address || '', // Populate the value from userData
              input_status: userData?.Status || '', // Populate the value from userData
              input_role: userData?.Rolename || '', // Populate the value from userData
              input_language: userData?.Language || '',
            }}
            validationSchema={validationSchema}
            onSubmit={handleSubmit}
          >
            {({ errors, touched, handleSubmit }) => (
              <form noValidate onSubmit={handleSubmit}>
                <Grid container spacing={2}>
                  <Grid item xs={3}>
                    <Field
                      name="fullName"
                      as={TextField}
                      label={t('fullNameLabel')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                    />
                    {touched.fullName && errors.fullName && (
                      <FormHelperText error>{errors.fullName}</FormHelperText>
                    )}
                  </Grid>
                  <Grid item xs={3}>
                    <Field
                      name="email"
                      as={TextField}
                      label={t('emailLabel')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                      InputProps={{
					   readOnly: true,
					 }}
					   inputProps={{
					      style: {
					      backgroundColor: 'rgb(208 208 214)',
					      },
					 }}
                    />
                    {touched.email && errors.email && (
                      <FormHelperText error>{errors.email}</FormHelperText>
                    )}
                  </Grid>
                  <Grid item xs={3}>
                    <Field
                      name="mobile"
                      as={TextField}
                      label={t('mobileNumberLabel')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                    />
                    {touched.mobile && errors.mobile && (
                      <FormHelperText error>{errors.mobile}</FormHelperText>
                    )}
                  </Grid>
                  <Grid item xs={3}>
                    <Field
                      name="address"
                      as={TextField}
                      label={t('addressLabel')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                    />
                    {touched.address && errors.address && (
                      <FormHelperText error>{errors.address}</FormHelperText>
                    )}
                  </Grid>
                {data && (
                    <>
                      <Grid item xs={3}>
                        <FormControl fullWidth variant="outlined" margin="normal">
                          <InputLabel id="role-label">{t('roleNameLabel')}</InputLabel>
                          <Field
                            name="input_role"
                            as={Select}
                            labelId="role-label"
                            label="Role"
                          >
                           <MenuItem value="">{t('Set Role')}</MenuItem>
                            {data.RoleData.Fields1
                              ? data.RoleData.Fields1.map((role) => (
                                <MenuItem key={role.Id} value={role.Id}>
                                  {role.Name}
                                </MenuItem>
                              ))
                              : <MenuItem disabled>No data available</MenuItem>
                            }
                          </Field>
                        </FormControl>
                        {touched.input_role && errors.input_role && (
                          <FormHelperText error>{errors.input_role}</FormHelperText>
                        )}
                      </Grid>
                    </>
                  )}
                <Grid item xs={3}>
                    <FormControl fullWidth variant="outlined" margin="normal">
                      <InputLabel id="language-label">{t('language')}</InputLabel>
                      <Field
                        name="input_language"
                        as={Select}
                        labelId="language-label"
                        label="Language"
                      >
                        <MenuItem value="">{t('selectLanguage')}</MenuItem>
                        {hardcodedLanguages.map((language) => (
                          <MenuItem key={language.id} value={language.id}>
                            {language.name}
                          </MenuItem>
                        ))}
                      </Field>
                      {touched.input_language && errors.input_language && (
                        <FormHelperText error>{errors.input_language}</FormHelperText>
                      )}
                    </FormControl>
                  </Grid>
                  <Grid item xs={3}>
                    <FormControl fullWidth variant="outlined" margin="normal">
                      <InputLabel id="status-label">{t('statusLabel')}</InputLabel>
                      <Field
                        name="input_status"
                        as={Select}
                        labelId="status-label"
                        label="Status"
                      >
                        <MenuItem value="">{t('selectStatus')}</MenuItem>
                        {hardcodedStatuses.map((status) => (
                          <MenuItem key={status.id} value={status.id}>
                            {status.name}
                          </MenuItem>
                        ))}
                      </Field>
                      {touched.input_status && errors.input_status && (
                        <FormHelperText error>{errors.input_status}</FormHelperText>
                      )}
                    </FormControl>
                  </Grid>
                </Grid>
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
                  <Button variant="contained" color="primary" style={{ backgroundColor: '#1478c8' }} type="submit">
                    {t('submit')}
                  </Button>
                  <Box sx={{ mx: 2 }} />
                  <Button variant="contained" color="secondary" style={{ backgroundColor: 'rgb(13 32 97)' }} onClick={() => navigate('/SysUser/SearchsysUser')}>
                    {t('back')}
                  </Button>
                </Box>
              </form>
            )}
          </Formik>
        </Card>
      </MainCard>
      <Dialog open={openPopup} onClose={handleDialogClose}>
        <DialogTitle>{t('superEsbAdminAlert')}</DialogTitle>
        <DialogContent>
          {popupMessage}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDialogClose} color="primary">
            Close
          </Button>
        </DialogActions>
      </Dialog>

    </>
  );
};

export default TablerIcons;