import React, { useState, useEffect } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
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
import { CREATE_CONSUMERS_URL, SEARCH_CONSUMERS_URL } from '../../../UrlPath.js';
import { fontWeight } from '@mui/system';
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

const CreateConsumers = () => {
  const navigate = useNavigate();
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
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
      <Typography color="text.primary">{t('createConsumer')}</Typography>
    </Breadcrumbs>
  ];

  // --------For status dropdown------
  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

  const [formData, setFormData] = useState({
    name: '',
    email: '',
    consumercode: '',
    consumerservices: '',
    input_status: '',
    consumeraddress: '',
  });

  // Function to parse JSON and check if it's valid
const parseJSON = (jsonString) => {
  try {
    return JSON.parse(jsonString);
  } catch (error) {
    return null;
  }
};

// Function to validate if the JSON object matches the specified structure
const isValidConsumerServices = (object) => {
  return (
    object &&
    object.services_list &&
    Array.isArray(object.services_list) &&
    object.services_list.every((service) =>
      service.service_name && typeof service.service_name === 'string' &&
      service.service_url && typeof service.service_url === 'string'
    )
  );
};

  const handleSubmit = async (values, { resetForm }) => {
    try {

       // Check if consumerservices is a valid JSON object
    const consumerservicesObject = parseJSON(values.consumerservices);

    if (!consumerservicesObject || !isValidConsumerServices(consumerservicesObject)) {
      openPopupWithMessage(t('validJsonObject'));
      return;
    }
    
      console.log('Sending data to backend:', values);
      const response = await axios.post('/Consumers/CreateConsumers', values);
      if (response.data.success) {
        openPopupWithMessage('Consumers created successfully');
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
  const validationSchema = Yup.object().shape({
    name: Yup.string()
      .matches(/^[A-Za-z\s]+$/, t('validationConsumerName1'))
      .max(30, t('validationConsumerName2'))
      .required(t('validationConsumerName3')),
    email: Yup.string()
      .max(30, t('validationEmail1'))
      .email(t('validationEmail2'))
      .required(t('validationEmail3')),
    consumercode: Yup.string()
      .max(30, t('validationConsumerCode1'))
      .required(t('validationConsumerCode2')),
    consumerservices: Yup.string().required(t('validationConsumerServices')),
    input_status: Yup.string().required(t('validationStatus')),
    consumeraddress: Yup.string()
      .max(30, t('validationConsumerAddress1'))
      .required(t('validationConsumerAddress2')),

  });

  //-----------For Popup Messages--------//

  const openPopupWithMessage = (message) => {
    setPopupMessage(message);
    setOpenPopup(true);
  };

  const structDefinition = `type ConsumerServices struct {
    ServicesList []struct {
      ServiceName   string \`json:"service_name"\`
      ServiceURL    string \`json:"service_url"\`
    } \`json:"services_list"\`
  }`;

  const SamplestructDefinition = `{
    "services_list":[
       {
          "service_name":"Get Customer Bill",
          "service_url":"/getbill"
       },
       {
          "service_name":"Get Status of Transaction",
          "service_url":"/transactionstatus"
       }
    ]
 }`;

  return (
    <>
      <MainCard title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('createConsumer')}</span>} breadcrumbs={searchUserBreadcrumbs}>
        <Card sx={{ overflow: 'hidden' }}>
          <Formik
            initialValues={{
              name: '',
              email: '',
              consumercode: '',
              consumerservices: '',
              input_status: '',
              consumeraddress: '',
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
                      label={t('consumerName')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
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
                    <Field
                      name="consumeraddress"
                      as={TextField}
                      label={t('consumerDomain')}
                      variant="outlined"
                      fullWidth
                      margin="normal"
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
                    />
                    {touched.consumercode && errors.consumercode && (
                      <FormHelperText error>{errors.consumercode}</FormHelperText>
                    )}
                  </Grid>


                  {/* <Grid item xs={3}>
                    <Field
                      name="endpointsconfig"
                      as={TextField}
                      label="End Point Configuration"
                      variant="outlined"
                      fullWidth
                      margin="normal"
                    />
                    {touched.address && errors.address && (
                      <FormHelperText error>{errors.address}</FormHelperText>
                    )}
                  </Grid> */}

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
                      />

                    </FormControl>

                    {touched.consumerservices && errors.consumerservices && (
                      <FormHelperText error>{errors.consumerservices}</FormHelperText>
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
                <Grid container spacing={2}>
                  <Grid item xs={6}>
                    <br />
                    <span style={{ fontWeight: 'bold' }}>{t('JsonStruct')}</span>
                    <br />
                    <pre>
                      <code>{structDefinition}</code>
                    </pre>

                  </Grid>
                  <Grid item xs={6}>
                    <br />
                    <span style={{ fontWeight: 'bold' }}>{t('JsonForReference')}</span>
                    <br />
                    <pre>
                      <code>{SamplestructDefinition}</code>
                    </pre>

                  </Grid>
                </Grid>
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
                  <Button variant="contained" color="primary" style={{ backgroundColor: '#1478c8' }} type="submit">
                    {t('submit')}
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

export default CreateConsumers;