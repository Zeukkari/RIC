<?php

/**
 * The plugin bootstrap file
 *
 * This file is read by WordPress to generate the plugin information in the plugin
 * admin area. This file also includes all of the dependencies used by the plugin,
 * registers the activation and deactivation functions, and defines a function
 * that starts the plugin.
 *
 * @link              https://phz.fi/
 * @since             0.0.0
 * @package           Ric_Wp
 *
 * @wordpress-plugin
 * Plugin Name:       RIC WordPress Plugin
 * Plugin URI:        https://github.com/phzfi/RIC
 * Description:       This is the Responsive Image Cache integration plugin for WordPress. It replaces WP media URL's with RIC URL's.
 * Version:           0.0.1
 * Author:            Nicholas Saarela
 * Author URI:        https://phz.fi/
 * License:           GPL-2.0+
 * License URI:       http://www.gnu.org/licenses/gpl-2.0.txt
 * Text Domain:       ric-wp
 * Domain Path:       /languages
 */

// If this file is called directly, abort.
if ( ! defined( 'WPINC' ) ) {
	die;
}

/**
 * The code that runs during plugin activation.
 * This action is documented in includes/class-ric-wp-activator.php
 */
function activate_ric_wp() {
	require_once plugin_dir_path( __FILE__ ) . 'includes/class-ric-wp-activator.php';
	Ric_Wp_Activator::activate();
}

/**
 * The code that runs during plugin deactivation.
 * This action is documented in includes/class-ric-wp-deactivator.php
 */
function deactivate_ric_wp() {
	require_once plugin_dir_path( __FILE__ ) . 'includes/class-ric-wp-deactivator.php';
	Ric_Wp_Deactivator::deactivate();
}

register_activation_hook( __FILE__, 'activate_ric_wp' );
register_deactivation_hook( __FILE__, 'deactivate_ric_wp' );

/**
 * The core plugin class that is used to define internationalization,
 * admin-specific hooks, and public-facing site hooks.
 */
require plugin_dir_path( __FILE__ ) . 'includes/class-ric-wp.php';

/**
 * Begins execution of the plugin.
 *
 * Since everything within the plugin is registered via hooks,
 * then kicking off the plugin from this point in the file does
 * not affect the page life cycle.
 *
 * @since    1.0.0
 */
function run_ric_wp() {

	$plugin = new Ric_Wp();
	$plugin->run();

}
run_ric_wp();
