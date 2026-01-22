package measurement

import "testing"

func BenchmarkCompare(b *testing.B) {
	m1 := Measurement{
		Type: TypeK8s,
		Subtypes: []Subtype{
			{
				Name: "cluster",
				Data: map[string]Reading{
					"version": Str("1.28.0"),
					"nodes":   Int(3),
					"ready":   Bool(true),
				},
			},
			{
				Name: "pod",
				Data: map[string]Reading{
					"count": Int(100),
					"ready": Int(95),
				},
			},
		},
	}

	m2 := Measurement{
		Type: TypeK8s,
		Subtypes: []Subtype{
			{
				Name: "cluster",
				Data: map[string]Reading{
					"version": Str("1.29.0"),
					"nodes":   Int(5),
					"ready":   Bool(true),
				},
			},
			{
				Name: "pod",
				Data: map[string]Reading{
					"count": Int(150),
					"ready": Int(140),
				},
			},
			{
				Name: "service",
				Data: map[string]Reading{
					"count": Int(50),
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(m1, m2)
	}
}

func BenchmarkCompare_NoChanges(b *testing.B) {
	m := Measurement{
		Type: TypeK8s,
		Subtypes: []Subtype{
			{
				Name: "cluster",
				Data: map[string]Reading{
					"version": Str("1.28.0"),
					"nodes":   Int(3),
					"ready":   Bool(true),
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(m, m)
	}
}

func BenchmarkCompare_ManySubtypes(b *testing.B) {
	// Create measurements with many subtypes
	subtypes1 := make([]Subtype, 50)
	subtypes2 := make([]Subtype, 50)

	for i := 0; i < 50; i++ {
		subtypes1[i] = Subtype{
			Name: "subtype" + string(rune(i)),
			Data: map[string]Reading{
				"value1": Int(i),
				"value2": Str("test"),
				"value3": Bool(i%2 == 0),
			},
		}
		subtypes2[i] = Subtype{
			Name: "subtype" + string(rune(i)),
			Data: map[string]Reading{
				"value1": Int(i + 1), // Changed
				"value2": Str("test"),
				"value3": Bool(i%2 == 0),
			},
		}
	}

	m1 := Measurement{Type: TypeK8s, Subtypes: subtypes1}
	m2 := Measurement{Type: TypeK8s, Subtypes: subtypes2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(m1, m2)
	}
}
