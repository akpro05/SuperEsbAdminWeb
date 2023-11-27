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
import { CREATE_SUBSCRIBER_URL, SEARCH_SUBSCRIBER_URL } from '../../../UrlPath.js';

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

import Breadcrumbs from '@mui/material/Breadcrumbs';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';

const UpdateProducer = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [userData, setUserData] = useState(null);
  const { t, i18n } = useTranslation();

  const searchUserBreadcrumbs = [
    <Breadcrumbs aria-label="breadcrumb">
      <Link underline="hover" color="inherit" href="/Dashboard">
        <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
        {t('Home')}
      </Link>
      <Link underline="hover" color="inherit" href="#">
      {t('Producer Management')}
      </Link>
      <Typography color="text.primary">{t('updateProducers')}</Typography>
    </Breadcrumbs>
  ];

  

  // --------For status dropdown------
  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

  const initialValues = {
    name: '',
    email: '',
    input_status: '',

  };

  const validationSchema = Yup.object().shape({
  name: Yup.string()
    .matches(/^[A-Za-z\s]+$/, t('validationFullName1'))
    .max(30, t('validationFullName2'))
    .required(t('validationFullName3')),

  email: Yup.string()
    .max(30, t('validationEmail1'))
    .email(t('validationEmail2'))
    .required(t('validationEmail3')),
  input_status: Yup.string().required(t('validationStatus')),
});

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await axios.get(`/Producers/UpdateProducers/${id}`);
        setUserData(response.data);
        setLoading(false);
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
    console.log('Form values to be sent to the backend:', values);
    try {
      const response = await axios.post(`/Producers/UpdateProducers/${id}`, values);
      if (response.data.message) {
        setPopupMessage(response.data.message);
        setOpenPopup(true); // Open the popup dialog

        // setValues({
        //   name: '',
        //   email: '',
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
      //   name: '',
      //   email: '',
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
    navigate('/Producers/SearchProducers'); // Redirect to /SysUser/SearchsysUser
  };

  return (
    <>
      <MainCard    title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('updateProducers')}</span>} breadcrumbs={searchUserBreadcrumbs}>
        <Card sx={{ overflow: 'hidden' }}>
          <Formik
            initialValues={{
              name: userData?.FullName || '', // Populate the value from userData
              email: userData?.Email || '', // Populate the value from userData
              input_status: userData?.Status || '', // Populate the value from userData
            }}
            validationSchema={validationSchema}
            onSubmit={handleSubmit}
          >
            {({ errors, touched, handleSubmit }) => (
              <form noValidate onSubmit={handleSubmit}>
                <Grid container spacing={2}>
                  <Grid item xs={3}>
                    <Field
                      name="name"
                      as={TextField}
                      label={t('producerName')}
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
                    {touched.name && errors.name && (
                      <FormHelperText error>{errors.name}</FormHelperText>
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
                    />
                    {touched.email && errors.email && (
                      <FormHelperText error>{errors.email}</FormHelperText>
                    )}
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

                    </FormControl>

                    {touched.input_status && errors.input_status && (
                      <FormHelperText error>{errors.input_status}</FormHelperText>
                    )}
                  </Grid>
                </Grid>
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
                  <Button variant="contained" color="primary" style={{ backgroundColor: '#1478c8' }} type="submit">
                    {t('submit')}
                  </Button>
                  <Box sx={{ mx: 2 }} />
                  <Button variant="contained" color="secondary" style={{ backgroundColor: 'rgb(13 32 97)' }} onClick={() => navigate('/Producers/SearchProducers')}>
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

export default UpdateProducer;