package ska

import (
	"crypto/aes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_generateKey(t *testing.T) {
	type args struct {
		input        string
		AESKeyLength AESKeyLength
	}
	type want struct {
		len int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "length 16",
			args: args{
				input:        "secret_key",
				AESKeyLength: Key16,
			},
			want: want{
				len: 16,
			},
		},
		{
			name: "length 24",
			args: args{
				input:        "secret_key",
				AESKeyLength: Key24,
			},
			want: want{
				len: 24,
			},
		},
		{
			name: "length 32",
			args: args{
				input:        "secret_key",
				AESKeyLength: Key32,
			},
			want: want{
				len: 32,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := generateKey(tt.args.input, tt.args.AESKeyLength)
			if len(got) != tt.want.len {
				t.Errorf("generateKey() = %d, want %d", len(got), tt.want.len)
			}
		})
	}
}

func Test_padData(t *testing.T) {
	type args struct {
		data      []byte
		blockSize int
	}
	type want struct {
		len int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "need to pad",
			args: args{
				data:      []byte("asd"),
				blockSize: aes.BlockSize,
			},
			want: want{
				len: aes.BlockSize,
			},
		},
		{
			name: "need to pad for two block size",
			args: args{
				data:      []byte("qwertasdfgzxcvbqw"),
				blockSize: aes.BlockSize,
			},
			want: want{
				len: 2 * aes.BlockSize,
			},
		},
		{
			name: "no need to pad",
			args: args{
				data:      []byte("qwertasdfgzxcvbq"),
				blockSize: aes.BlockSize,
			},
			want: want{
				len: aes.BlockSize,
			},
		},
		{
			name: "empty input",
			args: args{
				data:      []byte(""),
				blockSize: aes.BlockSize,
			},
			want: want{
				len: 0,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := padData(tt.args.data, tt.args.blockSize)
			if len(got)%tt.args.blockSize != 0 || len(got) != tt.want.len {
				t.Errorf("padData() = len %d, want %d", len(got), tt.want.len)
			}
		})
	}
}

func Test_unPadData(t *testing.T) {
	type args struct {
		data      []byte
		blockSize int
	}
	type want struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "base case",
			args: args{
				data:      []byte("asd"),
				blockSize: aes.BlockSize,
			},
			want: want{
				data: []byte("asd"),
			},
		},
		{
			name: "need to unpad from two block size",
			args: args{
				data:      []byte("qwertasdfgzxcvbqw"),
				blockSize: aes.BlockSize,
			},
			want: want{
				data: []byte("qwertasdfgzxcvbqw"),
			},
		},
		{
			name: "no need to pad",
			args: args{
				data:      []byte("qwertasdfgzxcvbq"),
				blockSize: aes.BlockSize,
			},
			want: want{
				data: []byte("qwertasdfgzxcvbq"),
			},
		},
		{
			name: "empty input",
			args: args{
				data:      []byte(""),
				blockSize: aes.BlockSize,
			},
			want: want{
				data: []byte(""),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			paddedData := padData(tt.args.data, tt.args.blockSize)

			got := unPadData(paddedData)
			if !reflect.DeepEqual(got, tt.want.data) {
				t.Errorf("unPadData() = %v, want %v", got, tt.want.data)
			}
		})
	}
}

func TestNewSKA(t *testing.T) {
	type args struct {
		userKey      string
		AESKeyLength AESKeyLength
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "base case 16",
			args: args{
				userKey:      "secret_key",
				AESKeyLength: Key16,
			},
		},
		{
			name: "base case 24",
			args: args{
				userKey:      "secret_key",
				AESKeyLength: Key24,
			},
		},
		{
			name: "base case 32",
			args: args{
				userKey:      "secret_key",
				AESKeyLength: Key32,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewSKA(tt.args.userKey, tt.args.AESKeyLength)
			assert.NotNil(t, got)
		})
	}
}

func TestSKA_SetAESKey(t *testing.T) {
	type args struct {
		userKey      string
		AESKeyLength AESKeyLength
	}
	type want struct {
		len int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "base 16",
			args: args{
				userKey:      "secret16",
				AESKeyLength: Key16,
			},
			want: want{
				len: 16,
			},
		},
		{
			name: "base 24",
			args: args{
				userKey:      "secret24",
				AESKeyLength: Key24,
			},
			want: want{
				len: 24,
			},
		},
		{
			name: "base 32",
			args: args{
				userKey:      "secret32",
				AESKeyLength: Key32,
			},
			want: want{
				len: 32,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := NewSKA("secret", Key16)
			current := s.keyAES
			s.SetAESKey(tt.args.userKey, tt.args.AESKeyLength)
			require.NotEqual(t, tt.args.userKey, s.keyAES)

			assert.NotEqual(t, current, s.keyAES)
			assert.Equal(t, len(s.keyAES), tt.want.len)
		})
	}
}

func TestSKA_Encrypt(t *testing.T) {
	type args struct {
		userKey string
		keyLen  AESKeyLength
		rawText []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "base 16",
			args: args{
				userKey: "secret",
				keyLen:  Key16,
				rawText: []byte("some text to encrypt/decrypt"),
			},
		},
		{
			name: "base 24",
			args: args{
				userKey: "secret",
				keyLen:  Key24,
				rawText: []byte("some text to encrypt/decrypt"),
			},
		},
		{
			name: "base 32",
			args: args{
				userKey: "secret",
				keyLen:  Key32,
				rawText: []byte("some text to encrypt/decrypt"),
			},
		},
		{
			name: "base empty",
			args: args{
				userKey: "secret",
				keyLen:  Key16,
				rawText: []byte(""),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := NewSKA(tt.args.userKey, tt.args.keyLen)
			encrypted, err := s.Encrypt(tt.args.rawText)
			require.NoError(t, err)
			assert.NotEqual(t, tt.args.rawText, encrypted)

			decrypted, err := s.Decrypt(encrypted)
			require.NoError(t, err)
			assert.Equal(t, tt.args.rawText, decrypted)
		})
	}
}
