package persistence

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/tucats/gopackages/app-cli/ui"
)

// ProfileDirectory is the name of the invisible directory that is created
// in the user's home directory to host configuration data
const ProfileDirectory = ".org.fernwood"

// ProfileFile is the name of the configuration file that contains the
// profiles.
var ProfileFile = "config.json"

// ProfileName is the name of the configuration being used. The default
// configuration is always named "default"
var ProfileName = "default"

// Configuration describes what is known about a configuration
type Configuration struct {
	Description string            `json:"description,omit"`
	ID          string            `json:"id,omit"`
	Items       map[string]string `json:"items"`
}

// CurrentConfiguration describes the current configuration that is active.
var CurrentConfiguration *Configuration

// explicitValues contains overridden default values
var explicitValues = Configuration{Description: "overridden defaults", Items: map[string]string{}}

// ProfileDirty is set to true when a key value is written or deleted, which
// tells us to rewrite the profile. If false, then no update is required.
var ProfileDirty = false

// Configurations is a map keyed by the configuration name for each
// configuration in the config file.
//
// @TOMCOLE This is currently a public symbol. Is this really needed?
// Check on the JSON code to determine if it can handle a private
// symbol at the top level...
var Configurations map[string]Configuration

// Load reads in the named profile, if it exists.
func Load(application string, name string) error {

	var c Configuration = Configuration{Description: "Default configuration", Items: map[string]string{}}

	CurrentConfiguration = &c
	Configurations = map[string]Configuration{"default": c}
	ProfileFile = application + ".json"

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
			ProfileDirty = true
		}
		ProfileName = name
		CurrentConfiguration = &c
	}

	return err
}

// Save the current configuration.
func Save() error {

	// So we even need to do anything?
	if !ProfileDirty {
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

	// Make sure every configuration has an id
	for n := range Configurations {
		c := Configurations[n]
		if c.ID == "" {
			c.ID = uuid.New().String()
			ui.Debug("Creating configuration \"%s\" with id %s", n, c.ID)
			Configurations[n] = c
		}
	}
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
		ProfileDirty = true
	}
	ProfileName = name
	CurrentConfiguration = &c
}

// Set puts a profile entry in the current Configuration structure
func Set(key string, value string) {

	explicitValues.Items[key] = value

	c := getCurrentConfiguration()
	c.Items[key] = value
	ProfileDirty = true
	ui.Debug("Setting profile key \"%s\" = \"%s\"", key, value)

}

// SetDefault puts a profile entry in the current Configuration structure. It is
// different than Set() in that it doesn't mark the value as dirty, so no need
// to update on account of this setting.
func SetDefault(key string, value string) {
	explicitValues.Items[key] = value
	ui.Debug("Setting default key \"%s\" = \"%s\"", key, value)
}

// Get gets a profile entry in the current configuration structure.
// If the key does not exist, an empty string is returned.
func Get(key string) string {

	// First, search the default values that be explicitly set

	v, found := explicitValues.Items[key]
	if !found {
		c := getCurrentConfiguration()
		v = c.Items[key]
	}
	ui.Debug("Reading profile key \"%s\" : \"%s\"", key, v)
	return v
}

// GetBool returns the boolean value of a profile string. If the string is
// "Y", "YES", "1", or "TRUE" then the value returns true.
func GetBool(key string) bool {
	s := strings.ToLower(Get(key))
	if s == "y" || s == "yes" || s == "true" || s == "t" || s == "1" {
		return true
	}
	return false
}

// Delete removes a key from the map entirely. Also removes if from the
// active defaults.
func Delete(key string) {
	c := getCurrentConfiguration()
	delete(c.Items, key)
	delete(explicitValues.Items, key)
	ProfileDirty = true
	ui.Debug("Deleting profile key \"%s\"", key)
}

// Keys returns the list of keys in the profile as an array
// of strings.
func Keys() []string {
	c := getCurrentConfiguration()
	result := []string{}
	for key := range c.Items {
		result = append(result, key)
	}
	return result
}

// Exists test to see if a key value exists or not
func Exists(key string) bool {

	_, exists := explicitValues.Items[key]
	if !exists {
		c := getCurrentConfiguration()
		_, exists = c.Items[key]
	}
	return exists
}

func getCurrentConfiguration() *Configuration {

	if CurrentConfiguration == nil {
		var c Configuration = Configuration{Description: "Default configuration", Items: map[string]string{}}
		CurrentConfiguration = &c
	}
	return CurrentConfiguration
}
