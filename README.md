# Intel® Generic Implementation of GraphBLAS* for Go*

This library is a binding for GraphBLAS, and more specifically for the SuiteSparse:GraphBLAS implementation of the GraphBLAS specification, for the Go programming language. It has been tested against version 8.0.0, 8.0.1, and 8.0.2 of SuiteSparse:GraphBLAS.

GraphBLAS is a specification for an API that defines standard building blocks for expressing graph algorithms in the language of linear algebra. It consists of data types for opaque representations of sparse matrices, vectors and scalars over the usual elementary types (Booloans, integers, floating point numbers), as well as user-defined element types. Operations on those data types include matrix-vector, vector-matrix, matrix-matrix multiplications, element-wise addition and multiplication, and so on. Apart from the usual integer and floating point addition and multiplication, client programs can use other arbitrary operators with these high-level operations, represented as monoids or semirings, to express a wide range of powerful algorithms.

SuiteSparse:GraphBLAS is a complete implementation of GraphBLAS in the C programming language. It is used heavily in production. For example, it is used as the underlying graph engine in RedisGraph, and as the built-in sparse matrix multiply in MATLAB. Several bindings exist for SuiteSparse:GraphBLAS for Python, Julia, and PostgreSQL.

Links with more information:
* [GraphBLAS](https://graphblas.org)
* [SuiteSparse:GraphBLAS](https://people.engr.tamu.edu/davis/GraphBLAS.html)

The forGraphBLASGo library is a binding that defines a Go API for GraphBLAS and the SuiteSparse:GraphBLAS extensions. It calls into the SuiteSparse:GraphBLAS C implementation. However, it strives to adhere to Go programming style as much as possible. Most prominently, it uses type parameters that have been introduced in Go 1.18 to make the various GraphBLAS object types generic, for added type safety. Other supported Go features include: using multiple return values instead of reference parameters, and Go-style error handling.

The library deviates from Go programming style when it comes to resource management: GraphBLAS objects tend to allocate a lot of memory, especially matrices, and they are allocated by the SuiteSparse:GraphBLAS implementation in C. Therefore, GraphBLAS objects are not managed by Go's garbage collector, but must be explicitly freed. While it is possible to register GraphBLAS objects with Go's garbage collector using `runtime.SetFinalizer`, it is advisable to rather call the respective `Free()` method in a `defer` statement.

The forGraphBLASGo library has been used in forLAGraphGo to express a number of well-known graph algorithms in Go. See [forLAGraphGo](https://github.com/intel/forLAGraphGo) to see forGraphBLASGo in action.

---

## Naming conventions

The GraphBLAS C Specification uses a naming convention where C identifiers available to GraphBLAS programs all start with the prefix `GrB_` (for example `GrB_Matrix`). SuiteSparse:GraphBLAS extends GraphBLAS with several new concepts and functions, and the corresponding identifiers start with `GxB_` (for example `GxB_Context`).

In forGraphBLASGo, the goal is to follow common Go-style naming conventions, and an obvious way to translate the C naming convention prefix `GrB_` is to use qualified identifiers with `GrB` as the package name (for example `GrB.Matrix`). This is indeed the path forGraphBLASGo takes.

Unfortunately, there is no easy way to distinguish between 'core' GraphBLAS identifiers and extensions. A separate `GxB` package is not feasible, because both core functionality and extensions share too many concepts internally.

forGraphBLASGo therefore instead opts to merely document extensions. Godoc comments will always mention whether something is a SuiteSparse:GraphBLAS extensions or even a forGraphBLASGo extension.

For some SuiteSparse:GraphBLAS extensions, this requires renaming some identifiers because they would otherways clash with core GraphBLAS names. For example, `GrB_Format` specifies the external format for matrix import and export functions, whereas `GxB_Format` specifies whether a matrix is laid out internally by row or by column. In forGraphBLASGo, the former remains `GrB.Format`, and the latter becomes `GrB.Layout`. These naming deviations are also mentioned in the godoc comments.

---

## Documentation

The introductory chapters (Chapters 1 and 2) of the [GraphBLAS C API Specification](https://graphblas.org/docs/GraphBLAS_API_C_v2.0.0.pdf) largely also apply to forGraphBLASGo. See the [Go API documentation](https://pkg.go.dev/github.com/intel/forGraphBLASGo) for how the concepts translate to forGraphBLASGo.

Except as otherwise noted, the API documention for forGraphBLASGo is derived from the The GraphBLAS C API Specification, version 2.0.0, authored by Benjamin Brock, Aydın Buluç, Timothy Mattson, Scott McMillan, and José Moreira. That material is licensed under a Creative Commons Attribution 4.0 license (http://creativecommons.org/licenses/by/4.0/legalcode).

Large sections are also derived from the [User Guide for SuiteSparse:GraphBLAS](https://github.com/DrTimothyAldenDavis/GraphBLAS/blob/stable/Doc/GraphBLAS_UserGuide.pdf) authored by Timothy A. Davis, especially regarding the SuiteSparse:GraphBLAS extensions that are incorporated into forGraphBLASGo.

---

## Licensing forGraphBLASGo

forGraphBLASGo is licensed under the BSD 3-Clause License

---

## Licensing and supporting SuiteSparse:GraphBLAS

SuiteSparse:GraphBLAS is released primarily under the Apache-2.0 license, because of how the project is supported by many organizations (NVIDIA, Redis, MIT Lincoln Lab, Intel, IBM, and Julia Computing), primarily through gifts to the Texas A&M Foundation.  Because of this support, and to facilitate the wide-spread use of GraphBLAS, the decision was made to give this library a permissive open-source license (Apache-2.0).  Currently all source code required to create the C-callable library libgraphblas.so is licensed with Apache-2.0, and there are no plans to change this.

However, just because this code is free to use doesn't make it zero-cost to create.  If you are using GraphBLAS in a commercial closed-source product and are not supporting its development, please consider supporting this project to ensure that it will continue to be developed and enhanced in the future.

To support the development of GraphBLAS, contact the author (davis@tamu.edu) or the Texas A&M Foundation (True Brown, tbrown@txamfoundation.com; or Kevin McGinnis, kmcginnis@txamfoundation.com) for details.

SuiteSparse:GraphBLAS, is copyrighted by Timothy A. Davis, (c) 2017-2023, All Rights Reserved.  davis@tamu.edu.
