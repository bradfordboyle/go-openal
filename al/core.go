// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// C-level binding for OpenAL's "al" API.
//
// Please consider using the Go-level binding instead.
// See http://connect.creativelabs.com/openal/Documentation/OpenAL%201.1%20Specification.htm
// for details about OpenAL not described here.
//
// OpenAL types are (in principle) mapped to Go types as
// follows:
//
//	ALboolean	bool	(al.h says char, but Go's bool should be compatible)
//	ALchar		uint8	(although al.h suggests int8, Go's uint8 (aka byte) seems better)
//	ALbyte		int8	(al.h says char, implying that char is signed)
//	ALubyte		uint8	(al.h says unsigned char)
//	ALshort		int16
//	ALushort	uint16
//	ALint		int32
//	ALuint		uint32
//	ALsizei		int32	(although that's strange, it's what OpenAL wants)
//	ALenum		int32	(although that's strange, it's what OpenAL wants)
//	ALfloat		float32
//	ALdouble	float64
//	ALvoid		not applicable (but see below)
//
// We also stick to these (not mentioned explicitly in
// OpenAL):
//
//	ALvoid*		unsafe.Pointer (but never exported)
//	ALchar*		string
//
// Finally, in places where OpenAL expects pointers to
// C-style arrays, we use Go slices if appropriate:
//
//	ALboolean*	[]bool
//	ALvoid*		[]byte (see Buffer.SetData() for example)
//	ALint*		[]int32
//	ALuint*		[]uint32 []Source []Buffer
//	ALfloat*	[]float32
//	ALdouble*	[]float64
//
// Overall, the correspondence of types hopefully feels
// natural enough. Note that many of these types do not
// actually occur in the API.
//
// The names of OpenAL constants follow the established
// Go conventions: instead of AL_FORMAT_MONO16 we use
// FormatMono16 for example.
//
// Conversion to Go's camel case notation does however
// lead to name clashes between constants and functions.
// For example, AL_DISTANCE_MODEL becomes DistanceModel
// which collides with the OpenAL function of the same
// name used to set the current distance model. We have
// to rename either the constant or the function, and
// since the function name seems to be at fault (it's a
// setter but doesn't make that obvious), we rename the
// function.
//
// In fact, we renamed plenty of functions, not just the
// ones where collisions with constants were the driving
// force. For example, instead of the Sourcef/GetSourcef
// abomination, we use Getf/Setf methods on a Source type.
// Everything should still be easily recognizable for
// OpenAL hackers, but this structure is a lot more
// sensible (and reveals that the OpenAL API is actually
// not such a bad design).
//
// There are a few cases where constants would collide
// with the names of types we introduced here. Since the
// types serve a much more important function, we renamed
// the constants in those cases. For example AL_BUFFER
// would collide with the type Buffer so it's name is now
// Buffer_ instead. Not pretty, but in many cases you
// don't need the constants anyway as the functionality
// they represent is probably available through one of
// the convenience functions we introduced as well. For
// example consider the task of attaching a buffer to a
// source. In C, you'd say alSourcei(sid, AL_BUFFER, bid).
// In Go, you can say sid.Seti(Buffer_, bid) as well, but
// you probably want to say sid.SetBuffer(bid) instead.
//
// TODO: Decide on the final API design; the current state
// has only specialized methods, none of the generic ones
// anymore; it exposes everything (except stuff we can't
// do) but I am not sure whether this is the right API for
// the level we operate on. Not yet anyway. Anyone?
package al

/*
#include <stdlib.h>
#include <AL/al.h>
#include "wrapper.h"
*/
import "C"
import "unsafe"

// General purpose constants. None can be used with SetDistanceModel()
// to disable distance attenuation. None can be used with Source.SetBuffer()
// to clear a Source of buffers.
const (
	None = 0;
	alFalse = 0;
	alTrue = 1;
)

// GetInteger() queries.
const (
	alDistanceModel = 0xD000;
)

// GetFloat() queries.
const (
	alDopplerFactor = 0xC000;
	alDopplerVelocity = 0xC001;
	alSpeedOfSound = 0xC003;
)

// GetString() queries.
const (
	alVendor = 0xB001;
	alVersion = 0xB002;
	alRenderer = 0xB003;
	alExtensions = 0xB004;
)

// Shared Source/Listener properties.
const (
	alPosition = 0x1004;
	alVelocity = 0x1006;
	alGain = 0x100A;
)

// Results from Source.State() query.
const (
	Initial = 0x1011;
	Playing = 0x1012;
	Paused = 0x1013;
	Stopped = 0x1014;
)

// Results from Source.Type() query.
const (
	Static = 0x1028;
	Streaming = 0x1029;
	Undetermined = 0x1030;
)

// TODO: Source properties.
// Regardless of what your al.h header may claim, Pitch
// only applies to Sources, not to Listeners. And I got
// that from Chris Robinson himself.
const (
	alSourceRelative = 0x202;
	alConeInnerAngle = 0x1001;
	alConeOuterAngle = 0x1002;
	alPitch = 0x1003;
	alDirection = 0x1005;
	alLooping = 0x1007;
	alBuffer = 0x1009;
	alMinGain = 0x100D;
	alMaxGain = 0x100E;
	alReferenceDistance = 0x1020;
	alRolloffFactor = 0x1021;
	alConeOuterGain = 0x1022;
	alMaxDistance = 0x1023;
	alSecOffset = 0x1024;
	alSampleOffset = 0x1025;
	alByteOffset = 0x1026;
)

func GetString(param int32) string {
	return C.GoString(C.walGetString(C.ALenum(param)));
}

func getBoolean(param int32) bool {
	return C.alGetBoolean(C.ALenum(param)) != alFalse;
}

func getInteger(param int32) int32 {
	return int32(C.alGetInteger(C.ALenum(param)));
}

func getFloat(param int32) float32 {
	return float32(C.alGetFloat(C.ALenum(param)));
}

func getDouble(param int32) float64 {
	return float64(C.alGetDouble(C.ALenum(param)));
}

// Renamed, was GetBooleanv.
func getBooleans(param int32, data []bool) {
	C.walGetBooleanv(C.ALenum(param), unsafe.Pointer(&data[0]));
}

// Renamed, was GetIntegerv.
func getIntegers(param int32, data []int32) {
	C.walGetIntegerv(C.ALenum(param), unsafe.Pointer(&data[0]));
}

// Renamed, was GetFloatv.
func getFloats(param int32, data []float32) {
	C.walGetFloatv(C.ALenum(param), unsafe.Pointer(&data[0]));
}

// Renamed, was GetDoublev.
func getDoubles(param int32, data []float64) {
	C.walGetDoublev(C.ALenum(param), unsafe.Pointer(&data[0]));
}

// Error codes from GetError()/for GetString().
const (
	NoError = alFalse;
	InvalidName = 0xA001;
	InvalidEnum = 0xA002;
	InvalidValue = 0xA003;
	InvalidOperation = 0xA004;
)

// GetError() returns the most recent error generated
// in the AL state machine.
func GetError() uint32 {
	return uint32(C.alGetError());
}

// Renamed, was DopplerFactor.
func SetDopplerFactor (value float32) {
	C.alDopplerFactor(C.ALfloat(value));
}

// Renamed, was DopplerVelocity.
func SetDopplerVelocity (value float32) {
	C.alDopplerVelocity(C.ALfloat(value));
}

// Renamed, was SpeedOfSound.
func SetSpeedOfSound (value float32) {
	C.alSpeedOfSound(C.ALfloat(value));
}

// Distance models for SetDistanceModel() and GetDistanceModel().
const (
	InverseDistance = 0xD001;
	InverseDistanceClamped = 0xD002;
	LinearDistance = 0xD003;
	LinearDistanceClamped = 0xD004;
	ExponentDistance = 0xD005;
	ExponentDistanceClamped = 0xD006;
)

// SetDistanceModel() changes the current distance model.
// Pass "None" to disable distance attenuation.
// Renamed, was DistanceModel.
func SetDistanceModel(model int32) {
	C.alDistanceModel(C.ALenum(model));
}

///// Source /////////////////////////////////////////////////////////

// Sources represent sound emitters in 3d space.
type Source uint32;

// NewSources() creates n sources.
// Renamed, was GenSources.
func NewSources(n int) (sources []Source) {
	sources = make([]Source, n);
	C.walGenSources(C.ALsizei(n), unsafe.Pointer(&sources[0]));
	return;
}

// DeleteSources() deletes the given sources.
func DeleteSources(sources []Source) {
	n := len(sources);
	C.walDeleteSources(C.ALsizei(n), unsafe.Pointer(&sources[0]));
}

// Renamed, was SourcePlayv.
func PlaySources(sources []Source) {
	C.walSourcePlayv(C.ALsizei(len(sources)), unsafe.Pointer(&sources[0]));
}

// Renamed, was SourceStopv.
func StopSources(sources []Source) {
	C.walSourceStopv(C.ALsizei(len(sources)), unsafe.Pointer(&sources[0]));
}

// Renamed, was SourceRewindv.
func RewindSources(sources []Source) {
	C.walSourceRewindv(C.ALsizei(len(sources)), unsafe.Pointer(&sources[0]));
}

// Renamed, was SourcePausev.
func PauseSources(sources []Source) {
	C.walSourcePausev(C.ALsizei(len(sources)), unsafe.Pointer(&sources[0]));
}

// Renamed, was Sourcef.
func (self Source) setf(param int32, value float32) {
	C.alSourcef(C.ALuint(self), C.ALenum(param), C.ALfloat(value));
}

// Renamed, was Source3f.
func (self Source) set3f(param int32, value1, value2, value3 float32) {
	C.alSource3f(C.ALuint(self), C.ALenum(param), C.ALfloat(value1), C.ALfloat(value2), C.ALfloat(value3));
}

// Renamed, was Sourcefv.
func (self Source) setfv(param int32, values []float32) {
	C.walSourcefv(C.ALuint(self), C.ALenum(param), unsafe.Pointer(&values[0]));
}

// Renamed, was Sourcei.
func (self Source) seti(param int32, value int32) {
	C.alSourcei(C.ALuint(self), C.ALenum(param), C.ALint(value));
}

// Renamed, was Source3i.
func (self Source) set3i(param int32, value1, value2, value3 int32) {
	C.alSource3i(C.ALuint(self), C.ALenum(param), C.ALint(value1), C.ALint(value2), C.ALint(value3));
}

// Renamed, was Sourceiv.
func (self Source) setiv(param int32, values []int32) {
	C.walSourceiv(C.ALuint(self), C.ALenum(param), unsafe.Pointer(&values[0]));
}

// Renamed, was GetSourcef.
func (self Source) getf(param int32) float32 {
	return float32(C.walGetSourcef(C.ALuint(self), C.ALenum(param)));
}

// Renamed, was GetSource3f.
func (self Source) get3f(param int32) (value1, value2, value3 float32) {
	var v1, v2, v3 float32;
	C.walGetSource3f(C.ALuint(self), C.ALenum(param), unsafe.Pointer(&v1),
		unsafe.Pointer(&v2), unsafe.Pointer(&v3));
	value1, value2, value3 = v1, v2, v3;
	return;
}

// Renamed, was GetSourcefv.
func (self Source) getfv(param int32, values []float32) {
	C.walGetSourcefv(C.ALuint(self), C.ALenum(param), unsafe.Pointer(&values[0]));
}

// Renamed, was GetSourcei.
func (self Source) geti(param int32) int32 {
	return int32(C.walGetSourcei(C.ALuint(self), C.ALenum(param)));
}

// Renamed, was GetSource3i.
func (self Source) get3i(param int32) (value1, value2, value3 int32) {
	var v1, v2, v3 int32;
	C.walGetSource3i(C.ALuint(self), C.ALenum(param), unsafe.Pointer(&v1),
		unsafe.Pointer(&v2), unsafe.Pointer(&v3));
	value1, value2, value3 = v1, v2, v3;
	return;
}

// Renamed, was GetSourceiv.
func (self Source) getiv(param int32, values []int32) {
	C.walGetSourceiv(C.ALuint(self), C.ALenum(param), unsafe.Pointer(&values[0]));
}

// Renamed, was SourcePlay.
func (self Source) Play() {
	C.alSourcePlay(C.ALuint(self));
}

// Renamed, was SourceStop.
func (self Source) Stop() {
	C.alSourceStop(C.ALuint(self));
}

// Renamed, was SourceRewind.
func (self Source) Rewind() {
	C.alSourceRewind(C.ALuint(self));
}

// Renamed, was SourcePause.
func (self Source) Pause() {
	C.alSourcePause(C.ALuint(self));
}

// Renamed, was SourceQueueBuffers.
func (self Source) QueueBuffers(buffers []Buffer) {
	C.walSourceQueueBuffers(C.ALuint(self), C.ALsizei(len(buffers)), unsafe.Pointer(&buffers[0]));
}

// Renamed, was SourceUnqueueBuffers.
func (self Source) UnqueueBuffers(buffers []Buffer) {
	C.walSourceUnqueueBuffers(C.ALuint(self), C.ALsizei(len(buffers)), unsafe.Pointer(&buffers[0]));
}

///// Convenience ////////////////////////////////////////////////////

// General

// NewSource() creates a single source.
// Convenience function, see NewSources().
func NewSource() Source {
	return Source(C.walGenSource());
}

// DeleteSource() deletes a single source.
// Convenience function, see DeleteSources().
func DeleteSource(source Source) {
	C.walDeleteSource(C.ALuint(source));
}

// Source

// Convenience method, see Source.QueueBuffers().
func (self Source) QueueBuffer(buffer Buffer) {
	C.walSourceQueueBuffer(C.ALuint(self), C.ALuint(buffer));
}

// Convenience method, see Source.QueueBuffers().
func (self Source) UnqueueBuffer() Buffer {
	return Buffer(C.walSourceUnqueueBuffer(C.ALuint(self)));
}

// Source queries.
// TODO: SourceType isn't documented as a query in the
// al.h header, but it is documented that way in
// the OpenAL 1.1 specification.
const (
	alSourceState = 0x1010;
	alBuffersQueued = 0x1015;
	alBuffersProcessed = 0x1016;
	alSourceType = 0x1027;
)

// Convenience method, see Source.Geti().
func (self Source) BuffersQueued() int32 {
	return self.geti(alBuffersQueued);
}

// Convenience method, see Source.Geti().
func (self Source) BuffersProcessed() int32 {
	return self.geti(alBuffersProcessed);
}

// Convenience method, see Source.Geti().
func (self Source) State() int32 {
	return self.geti(alSourceState);
}

// Convenience method, see Source.Geti().
func (self Source) Type() int32 {
	return self.geti(alSourceType);
}

// Convenience method, see Source.Getf().
func (self Source) GetGain() (gain float32) {
	return self.getf(alGain);
}

// Convenience method, see Source.Setf().
func (self Source) SetGain(gain float32) {
	self.setf(alGain, gain);
}

// Convenience method, see Source.Getf().
func (self Source) GetMinGain() (gain float32) {
	return self.getf(alMinGain);
}

// Convenience method, see Source.Setf().
func (self Source) SetMinGain(gain float32) {
	self.setf(alMinGain, gain);
}

// Convenience method, see Source.Getf().
func (self Source) GetMaxGain() (gain float32) {
	return self.getf(alMaxGain);
}

// Convenience method, see Source.Setf().
func (self Source) SetMaxGain(gain float32) {
	self.setf(alMaxGain, gain);
}

// Convenience method, see Source.Getf().
func (self Source) GetReferenceDistance() (distance float32) {
	return self.getf(alReferenceDistance);
}

// Convenience method, see Source.Setf().
func (self Source) SetReferenceDistance(distance float32) {
	self.setf(alReferenceDistance, distance);
}

// Convenience method, see Source.Getf().
func (self Source) GetMaxDistance() (distance float32) {
	return self.getf(alMaxDistance);
}

// Convenience method, see Source.Setf().
func (self Source) SetMaxDistance(distance float32) {
	self.setf(alMaxDistance, distance);
}

// Convenience method, see Source.Getf().
func (self Source) GetPitch() (gain float32) {
	return self.getf(alPitch);
}

// Convenience method, see Source.Setf().
func (self Source) SetPitch(gain float32) {
	self.setf(alPitch, gain);
}

// Convenience method, see Source.Getf().
func (self Source) GetRolloffFactor() (gain float32) {
	return self.getf(alRolloffFactor);
}

// Convenience method, see Source.Setf().
func (self Source) SetRolloffFactor(gain float32) {
	self.setf(alRolloffFactor, gain);
}

// Convenience method, see Source.Geti().
func (self Source) GetLooping() bool {
	return self.geti(alLooping) != alFalse;
}

var bool2al map[bool]int32 = map[bool]int32{true: alTrue, false: alFalse}

// Convenience method, see Source.Seti().
func (self Source) SetLooping(yes bool) {
	self.seti(alLooping, bool2al[yes]);
}

// Convenience method, see Source.Geti().
func (self Source) GetSourceRelative() bool {
	return self.geti(alSourceRelative) != alFalse;
}

// Convenience method, see Source.Seti().
func (self Source) SetSourceRelative(yes bool) {
	self.seti(alSourceRelative, bool2al[yes]);
}

// Convenience method, see Source.Setfv().
func (self Source) SetPosition(vector Vector) {
	self.setfv(alPosition, vector[0:]);
}

// Convenience method, see Source.Getfv().
func (self Source) GetPosition() Vector {
	v := Vector{};
	self.getfv(alPosition, v[0:]);
	return v;
}

// Convenience method, see Source.Setfv().
func (self Source) SetDirection(vector Vector) {
	self.setfv(alDirection, vector[0:]);
}

// Convenience method, see Source.Getfv().
func (self Source) GetDirection() Vector {
	v := Vector{};
	self.getfv(alDirection, v[0:]);
	return v;
}

// Convenience method, see Source.Setfv().
func (self Source) SetVelocity(vector Vector) {
	self.setfv(alVelocity, vector[0:]);
}

// Convenience method, see Source.Getfv().
func (self Source) GetVelocity() Vector {
	v := Vector{};
	self.getfv(alVelocity, v[0:]);
	return v;
}

// Convenience method, see Source.Getf().
func (self Source) GetOffsetSeconds() float32 {
	return self.getf(alSecOffset);
}

// Convenience method, see Source.Setf().
func (self Source) SetOffsetSeconds(offset float32) {
	self.setf(alSecOffset, offset);
}

// Convenience method, see Source.Geti().
func (self Source) GetOffsetSamples() int32 {
	return self.geti(alSampleOffset);
}

// Convenience method, see Source.Seti().
func (self Source) SetOffsetSamples(offset int32) {
	self.seti(alSampleOffset, offset);
}

// Convenience method, see Source.Geti().
func (self Source) GetOffsetBytes() int32 {
	return self.geti(alByteOffset);
}

// Convenience method, see Source.Seti().
func (self Source) SetOffsetBytes(offset int32) {
	self.seti(alByteOffset, offset);
}

// Convenience method, see Source.Getf().
func (self Source) GetInnerAngle() float32 {
	return self.getf(alConeInnerAngle);
}

// Convenience method, see Source.Setf().
func (self Source) SetInnerAngle(offset float32) {
	self.setf(alConeInnerAngle, offset);
}

// Convenience method, see Source.Getf().
func (self Source) GetOuterAngle() float32 {
	return self.getf(alConeOuterAngle);
}

// Convenience method, see Source.Setf().
func (self Source) SetOuterAngle(offset float32) {
	self.setf(alConeOuterAngle, offset);
}

// Convenience method, see Source.Getf().
func (self Source) GetOuterGain() float32 {
	return self.getf(alConeOuterGain);
}

// Convenience method, see Source.Setf().
func (self Source) SetOuterGain(offset float32) {
	self.setf(alConeOuterGain, offset);
}

// Convenience method, see Source.Geti().
func (self Source) SetBuffer(buffer Buffer) {
	self.seti(alBuffer, int32(buffer));
}

// Convenience method, see Source.Geti().
func (self Source) GetBuffer() (buffer Buffer) {
	return Buffer(self.geti(alBuffer));
}

///// Crap ///////////////////////////////////////////////////////////

// These functions are wrapped and should work fine, but they
// have no purpose: There are *no* capabilities in OpenAL 1.1
// which is the latest specification. So we removed from from
// the API for now, it's complicated enough without them.
//
//func Enable(capability int32) {
//	C.alEnable(C.ALenum(capability));
//}
//
//func Disable(capability int32) {
//	C.alDisable(C.ALenum(capability));
//}
//
//func IsEnabled(capability int32) bool {
//	return C.alIsEnabled(C.ALenum(capability)) != alFalse;
//}

// These constants are documented as "not yet exposed". We
// keep them here in case they ever become valid. They are
// buffer states.
//
//const (
//	Unused = 0x2010;
//	Pending = 0x2011;
//	Processed = 0x2012;
//)

// These functions would work fine, but they are not very
// useful since we have distinct Source and Buffer types.
// Leaving them out reduces API complexity, a good thing.
//
//func IsSource(id uint32) bool {
//	return C.alIsSource(C.ALuint(id)) != alFalse;
//}
//
//func IsBuffer(id uint32) bool {
//	return C.alIsBuffer(C.ALuint(id)) != alFalse;
//}
