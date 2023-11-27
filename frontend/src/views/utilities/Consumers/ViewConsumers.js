import React, { useState, useEffect } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import { Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton } from '@mui/material';

// Import the daterangepicker CSS
import axios from 'axios';
import { Formik, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';

// project imports
import MainCard from '../../../ui-component/cards/MainCard';
import HomeIcon from '@mui/icons-material/Home';
import { VIEW_CONSUMERS_URL } from '../../../UrlPath.js';
import Breadcrumbs from '@mui/material/Breadcrumbs';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';

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

const ViewConsumers = () => {
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
      {t('Consumer Management')}
      </Link>
      <Typography color="text.primary">{t('viewConsumer')}</Typography>
    </Breadcrumbs>
  ];
  // --------For status dropdown------
  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await axios.get(`/Consumers/ViewConsumers/${id}`);
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
  const ConsumerName = userData && userData.ConsumerName;
  const Email = userData && userData.Email;
  const ConsumerCode = userData && userData.ConsumerCode;	
  const Status = userData && userData.Status;
  const consumerservices = userData && userData.consumerservices;
  const consumeraddress = userData && userData.consumeraddress;


  //-----------For Popup Messages--------//

  const openPopupWithMessage = (message) => {
    setPopupMessage(message);
    setOpenPopup(true);
  };
  return (
    <>
      <MainCard  title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('viewConsumer')}</span>} breadcrumbs={searchUserBreadcrumbs}>
        <Card sx={{ overflow: 'hidden' }}>
          <Formik
            initialValues={{
              ConsumerName: '',
              email: '',
            ConsumerCode: '',
              input_status: '',
            consumerservices: '',
            consumeraddress: '',
            }}
          >
            {({ errors, touched, handleSubmit }) => (
              <form noValidate onSubmit={handleSubmit}>
                <Grid container spacing={2}>
                  <Grid item xs={3}>
                    <Field
                      name="name"
                      as={TextField}
                      label={t('consumerName')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                      value={ConsumerName}
                      InputProps={{
					      readOnly: true,
					   }}
					   inputProps={{
					      style: {
					        backgroundColor: 'rgb(208 208 214)',
					      },
					   }}
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
                      value={Email}
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
                      name="consumeraddress"
                      as={TextField}
                      label={t('consumerDomain')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                      value={consumeraddress}
                      InputProps={{
					      readOnly: true,
					   }}
					   inputProps={{
					      style: {
					        backgroundColor: 'rgb(208 208 214)',
					      },
					   }}
                    />
                    {touched.consumeraddress && errors.consumeraddress && (
                      <FormHelperText error>{errors.consumeraddress}</FormHelperText>
                    )}
                  </Grid>
                <Grid item xs={3}>
                    <Field
                      name="consumercode"
                      as={TextField}
                      label={t('consumerCode')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                      value={ConsumerCode}
                      InputProps={{
					      readOnly: true,
					   }}
					   inputProps={{
					      style: {
					        backgroundColor: 'rgb(208 208 214)',
					      },
					   }}
                    />
                    {touched.consumercode && errors.consumercode && (
                      <FormHelperText error>{errors.consumercode}</FormHelperText>
                    )}
                  </Grid>
                  
                  
                  
                 
                <Grid item xs={6}>
                    <FormControl fullWidth variant="outlined" margin="normal">
                      <InputLabel htmlFor="amf_metadata"></InputLabel>
                      <Field
                        as={TextField}
                        id="consumerservices"
                        name="consumerservices"
                        label={t('consumerServices')}
                        variant="outlined"
                        multiline
                        value={consumerservices}
                        InputProps={{
					      readOnly: true,
					   }}
					   inputProps={{
					      style: {
					        backgroundColor: 'rgb(208 208 214)',
					      },
					   }}
                      />
                       
                    </FormControl>

                    {touched.endpointsconfig && errors.endpointsconfig && (
                      <FormHelperText error>{errors.endpointsconfig}</FormHelperText>
                    )}
                  </Grid>
                  <Grid item xs={3}>
                    <Field
                      name="input_status"
                      as={TextField}
                      label={t('statusLabel')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
                      value={Status}
                      InputProps={{
					      readOnly: true,
					   }}
					   inputProps={{
					      style: {
					        backgroundColor: 'rgb(208 208 214)',
					      },
					   }}
                    />
                    {touched.address && errors.address && (
                      <FormHelperText error>{errors.address}</FormHelperText>
                    )}
                  </Grid>
                </Grid>
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
                  <Button variant="contained" color="primary" type="submit" style={{ backgroundColor: '#1478c8' }} onClick={() => navigate(`/Consumers/UpdateConsumers/${id}`)}>
                    {t('updateButton')}
                  </Button>
                  <Box sx={{ mx: 2 }} />
                  <Button variant="contained" color="secondary" style={{ backgroundColor: 'rgb(13 32 97)' }} onClick={() => navigate('/Consumers/SearchConsumers')}>
                    {t('back')}
                  </Button>
                </Box>
              </form>
            )}
          </Formik>
        </Card>
      </MainCard>
      <Dialog open={openPopup} onClose={() => setOpenPopup(false)}>
        <DialogTitle>SuperEsb Admin Alert</DialogTitle>
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

export default ViewConsumers;