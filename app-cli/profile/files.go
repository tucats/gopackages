// Package profile manages the persistent user profile used by the command
// application infrastructure. This includes automatically reading any
// profile in as part of startup, and of updating the profile as needed.
package profile

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// ProfileDirectory is the name of the invisible directory that is created
// in the user's home directory to host configuration data
var ProfileDirectory = ".org.fernwood"

// ProfileFile is the name of the configuration file that contains the
// profiles.
var ProfileFile = "config.json"

// ProfileName is the name of the configuration being used. The default
// configuration is always named "default"
var ProfileName = "default"

// Configuration describes what is known about a configuration
type Configuration struct {
	Description string            `json:"description,omit"`
	Items       map[string]string `json:"items"`
}

// CurrentConfiguration describes the current configuration that is active.
var CurrentConfiguration *Configuration

// explicitValues contains overridden default values
var explicitValues = Configuration{Description: "overridden defaults", Items: map[string]string{}}

// profileDirty is set to true when a key value is written or deleted, which
// tells us to rewrite the profile. If false, then no update is required.
var profileDirty = false

// Configurations is a map keyed by the configuration name for each
// configuration in the config file
var Configurations map[string]Configuration

// Load reads in the named profile, if it exists.
func Load(name string) error {

	var c Configuration = Configuration{Description: "Default configuration", Items: map[string]string{}}
	CurrentConfiguration = &c
	Configurations = map[string]Configuration{"default": c}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	var path strings.Builder
	path.WriteString(home)
	path.WriteRune(os.PathSeparator)
	path.WriteString(ProfileDirectory)
	path.WriteRune(os.PathSeparator)
	path.WriteString(ProfileFile)

	configFile, err := os.Open(path.String())
	if err != nil {
		return err
	}

	defer configFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(configFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into the config map which we defined above
	err = json.Unmarshal(byteValue, &Configurations)
	if err == nil {
		if name == "" {
			name = ProfileName
		}
		c, found := Configurations[name]
		if !found {
			c = Configuration{Description: "Default configuration", Items: map[string]string{}}
			Configurations[name] = c
		}
		ProfileName = name
		CurrentConfiguration = &c
	}

	return err
}

// Save the current configuration.
func Save() error {

	// So we even need to do anything?
	if !profileDirty {
		return nil
	}

	// Does the directory exist?
	var path strings.Builder
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path.WriteString(home)
	path.WriteRune(os.PathSeparator)
	path.WriteString(ProfileDirectory)

	if _, err := os.Stat(path.String()); os.IsNotExist(err) {
		os.MkdirAll(path.String(), os.ModePerm)
	}

	path.WriteRune(os.PathSeparator)
	path.WriteString(ProfileFile)

	byteBuffer, err := json.MarshalIndent(&Configurations, "", "  ")

	err = ioutil.WriteFile(path.String(), byteBuffer, os.ModePerm)
	return err
}

// UseProfile specifies the name of the profile to use, if other
// than the default.
func UseProfile(name string) {

	c, found := Configurations[name]
	if !found {
		c = Configuration{Description: name + " configuration", Items: map[string]string{}}
		Configurations[name] = c
		profileDirty = true
	}
	CurrentConfiguration = &c
}

// Set puts a profile entry in the current Configuration structure
func Set(key string, value string) {

	c := *CurrentConfiguration
	c.Items[key] = value
	profileDirty = true

}

// SetDefault puts a profile entry in the current Configuration structure. It is
// different than Set() in that it doesn't mark the value as dirty, so no need
// to update on account of this setting.
func SetDefault(key string, value string) {
	explicitValues.Items[key] = value
}

// Get gets a profile entry in the current configuration structure.
// If the key does not exist, an empty string is returned.
func Get(key string) string {

	// First, search the default values that be explicitly set

	v, found := explicitValues.Items[key]
	if !found {
		c := *CurrentConfiguration
		v = c.Items[key]
	}
	return v
}

// Delete removes a key from the map entirely. Also removes if from the
// active defaults.
func Delete(key string) {
	c := *CurrentConfiguration
	delete(c.Items, key)
	delete(explicitValues.Items, key)
	profileDirty = true
}

// Exists test to see if a key value exists or not
func Exists(key string) bool {

	_, exists := explicitValues.Items[key]
	if !exists {
		c := *CurrentConfiguration
		_, exists = c.Items[key]
	}
	return exists
}
