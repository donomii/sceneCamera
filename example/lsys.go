package main

import (
	"gitlab.com/donomii/wasm-lsystem/lsystem"
	"github.com/go-gl/mathgl/mgl32"
)

func calcLsys(CurrentScene *lsystem.Scene) ([]float32, []float32, []lsystem.HotSpot) {
	movMatrix := mgl32.Ident4()
	// movMatrix = Move(movMatrix, 1.0, 0.0, 0.0)

	verticesNative, colorsNative, hotspots := lsystem.Draw(CurrentScene,  lsystem.S(`
			Colour254,254,254 
			 s s s s s s s s HotSpot(1) Tetrahedron HotSpot(2) 
			deg30   [ HR
				s s
				[ s s HR Icosahedron ] TF TF TF TF 
				[ HR Tetrahedron ] Arrow  F  Arrow  F  Arrow  F  
				[ p p p s s s HR starburst ] Arrow  F  Arrow  F  Arrow  F 
				[ p p p s s HR leaf ] Arrow  F  Arrow  F  Arrow  F 	
				
				[ p p p s s s HR lineStar ] TF TF TF
				[ p p p s s HR Flower ] TF TF TF
				[ p p p s s HR Flower12 ] TF TF TF
				[ p p p s s HR Flower11 ] TF TF TF
				[ p p p s s HR Flower10 ] TF TF TF
				
				
			]
			
			p p p F P P P
			[ s s s s
			
				
				[ p p p S S S HR Square1 ] TF TF TF
				[ p p p S S S S S S HR Face ] TF TF TF
				[ p p p S S S HR Arrow ] TF TF TF
				[ p p p S HR Prism ] TF TF TF
				[ p p p S HR Prism1 ] TF TF TF
				[   s s HR p p p Circle ] TF TF TF

				
			]
			
			`), movMatrix, true)

	return verticesNative, colorsNative, hotspots
}
