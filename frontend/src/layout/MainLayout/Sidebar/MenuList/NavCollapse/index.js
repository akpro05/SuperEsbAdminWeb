import React from 'react';
import PropTypes from 'prop-types';
import { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { useLocation, useNavigate } from 'react-router';
import { useTranslation } from 'react-i18next';

// material-ui
import { useTheme } from '@mui/material/styles';
import { Collapse, List, ListItemButton, ListItemIcon, ListItemText, Typography } from '@mui/material';

// project imports
import NavItem from '../NavItem';

// assets
import FiberManualRecordIcon from '@mui/icons-material/FiberManualRecord';
import { IconChevronDown, IconChevronUp } from '@tabler/icons';

// ==============================|| SIDEBAR MENU LIST COLLAPSE ITEMS ||============================== //

const NavCollapse = ({ menu, level }) => {
  const theme = useTheme();
  const customization = useSelector((state) => state.customization);
  const navigate = useNavigate();
  const { t } = useTranslation();

  const [open, setOpen] = useState(false);
  const [selected, setSelected] = useState(null);

  const handleClick = () => {
    setOpen(!open);
    setSelected(!selected ? menu.id : null);
    if (menu?.id !== 'authentication') {
      navigate(menu.children[0]?.url);
    }
  };

  const { pathname } = useLocation();
  const checkOpenForParent = (child, id) => {
    child.forEach((item) => {
      if (item.url === pathname) {
        setOpen(true);
        setSelected(id);
      }
    });
  };

  // menu collapse for sub-levels
  useEffect(() => {
    setOpen(false);
    setSelected(null);



    if (menu.children) {
      menu.children.forEach((item) => {
        if (item.children?.length) {
          checkOpenForParent(item.children, menu.id);
        }
        if (item.url === pathname) {
          setSelected(menu.id);
          setOpen(true);
        }
      });
    }


    if (menu.title === 'System User Management' && (
      pathname.includes('/SysUser/CreatesysUser') ||
      pathname.match(/\/SysUser\/UpdatesysUser\/[a-f0-9-]+/) ||
      pathname.match(/\/SysUser\/ViewsysUser\/[a-f0-9-]+/)
    )) {
      setSelected(menu.id);
      setOpen(true);
    }

    if (menu.title === 'Producer Management' && (
      pathname.includes('/Producers/CreateProducers') ||
      pathname.match(/\/Producers\/UpdateProducers\/[a-f0-9-]+/) ||
      pathname.match(/\/Producers\/ViewProducers\/[a-f0-9-]+/)
    )) {
      setSelected(menu.id);
      setOpen(true);
    }

    if (menu.title === 'Consumer Management' && (
      pathname.includes('/Consumers/CreateConsumers') ||
      pathname.match(/\/Consumers\/UpdateConsumers\/[a-f0-9-]+/) ||
      pathname.match(/\/Consumers\/ViewConsumers\/[a-f0-9-]+/)
    )) {
      setSelected(menu.id);
      setOpen(true);
    }

    if (menu.title === 'System Configuration' && (
      pathname.includes('/ProducerToConsumer/CreateProducerToConsumer') ||
      pathname.match(/\/ProducerToConsumer\/UpdateProducerToConsumer\/[a-f0-9-]+/) ||
      pathname.match(/\/ProducerToConsumer\/ViewProducerToConsumer\/[a-f0-9-]+/)
    )) {
      setSelected(menu.id);
      setOpen(true);
    }


    if (menu.title === 'System Configuration' && (
      pathname.includes('/Role/CreateRole') ||
      pathname.match(/\/Role\/UpdateRole\/[a-f0-9-]+/) ||
      pathname.match(/\/Role\/ViewRole\/[a-f0-9-]+/)
    )) {
      setSelected(menu.id);
      setOpen(true);
    }


    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pathname, menu.children]);

  // menu collapse & item
  const menus = menu.children?.map((item) => {
    switch (item.type) {
      case 'collapse':
        return <NavCollapse key={item.id} menu={item} level={level + 1} />;
      case 'item':
        return <NavItem key={item.id} item={item} level={level + 1} />;
      default:
        return (
          <Typography key={item.id} variant="h6" color="error" align="center">
            Menu Items Error
          </Typography>
        );
    }
  });

  const Icon = menu.icon;
  const menuIcon = menu.icon ? (
    <Icon strokeWidth={1.5} size="1.3rem" style={{ marginTop: 'auto', marginBottom: 'auto' }} />
  ) : (
    <FiberManualRecordIcon
      sx={{
        width: selected === menu.id ? 8 : 6,
        height: selected === menu.id ? 8 : 6
      }}
      fontSize={level > 0 ? 'inherit' : 'medium'}
    />
  );

  return (
    <>
      <ListItemButton
        sx={{
          borderRadius: `${customization.borderRadius}px`,
          mb: 0.5,
          alignItems: 'flex-start',
          backgroundColor: level > 1 ? 'transparent !important' : 'inherit',
          py: level > 1 ? 1 : 1.25,
          pl: `${level * 24}px`
        }}
        selected={selected === menu.id}
        onClick={handleClick}
      >
        <ListItemIcon sx={{ my: 'auto', minWidth: !menu.icon ? 18 : 36 }}>{menuIcon}</ListItemIcon>
        <ListItemText
          primary={
            <Typography variant={selected === menu.id ? 'h5' : 'body1'} color="inherit" sx={{ my: 'auto' }}>
              {t(menu.title)}
            </Typography>
          }
          secondary={
            menu.caption && (
              <Typography variant="caption" sx={{ ...theme.typography.subMenuCaption }} display="block" gutterBottom>
                {t(menu.caption)}
              </Typography>
            )
          }
        />
        {open ? (
          <IconChevronUp stroke={1.5} size="1rem" style={{ marginTop: 'auto', marginBottom: 'auto' }} />
        ) : (
          <IconChevronDown stroke={1.5} size="1rem" style={{ marginTop: 'auto', marginBottom: 'auto' }} />
        )}
      </ListItemButton>
      <Collapse in={open} timeout="auto" unmountOnExit>
        <List
          component="div"
          disablePadding
          sx={{
            position: 'relative',
            '&:after': {
              content: "''",
              position: 'absolute',
              left: '32px',
              top: 0,
              height: '100%',
              width: '1px',
              opacity: 1,
              background: theme.palette.primary.light
            }
          }}
        >
          {menus}
        </List>
      </Collapse>
    </>
  );
};

NavCollapse.propTypes = {
  menu: PropTypes.object,
  level: PropTypes.number
};

export default NavCollapse;
